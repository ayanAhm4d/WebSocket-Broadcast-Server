package server

import (
	"broadcast/internal/config"
	"broadcast/pkg/websocket"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Client stuct represents the websocket client
type Client struct {
	cfg        *config.ClientConfig
	conn       *websocket.Connection
	mutex      sync.RWMutex
	done       chan struct{}
	msgChan    chan []byte
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewClient initializes a new websocket client
func NewClient(cfg *config.ClientConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		cfg:        cfg,
		done:       make(chan struct{}),
		msgChan:    make(chan []byte),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Connect establishes a connection to the WebSocket server and starts the client
func (c *Client) Connect() error {
	// Create WebSocket connection
	url := fmt.Sprintf("ws://%s/ws", c.cfg.ServerAddr)
	conn, _, err := ws.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.conn = websocket.NewConnection(conn)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("\nDisconnecting...")
		c.Disconnect()
	}()

	// Start client routines
	go c.readMessages()
	go c.processMessages()
	c.handleUserInput()

	return nil
}

// readMessages continuously reads messages from the WebSocket connection
func (c *Client) readMessages() {
	defer close(c.msgChan)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.conn.ReadMessage()
			if err != nil {
				if !ws.IsCloseError(err, ws.CloseNormalClosure, ws.CloseGoingAway) {
					log.Printf("Read error: %v", err)
				}
				c.Disconnect()
				return
			}
			c.msgChan <- msg
		}
	}
}

// processMessages handles incoming messages from the server
func (c *Client) processMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.msgChan:
			if !ok {
				return
			}
			c.displayMessage(msg)
		}
	}
}

// handleUserInput reads user input from stdin and sends it to the server
func (c *Client) handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return
		default:
			text := scanner.Text()
			if text == "" {
				fmt.Print("> ")
				continue
			}

			// Format message with username if provided
			var message string
			if c.cfg.Username != "" {
				message = fmt.Sprintf("[%s] %s", c.cfg.Username, text)
			} else {
				message = text
			}

			if err := c.sendMessage([]byte(message)); err != nil {
				log.Printf("Failed to send message: %v", err)
				c.Disconnect()
				return
			}
			fmt.Print("> ")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

// sendMessage sends a message to the server
func (c *Client) sendMessage(message []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	return c.conn.WriteMessage(message)
}

// displayMessage formats and displays a received message
func (c *Client) displayMessage(message []byte) {
	// Clear the current input line
	fmt.Printf("\r%s\n> ", string(message))
}

// Disconnect closes the client connection
func (c *Client) Disconnect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		c.cancelFunc()

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Send close message to server
		message := ws.FormatCloseMessage(ws.CloseNormalClosure, "")
		_ = c.conn.WriteMessage([]byte(message))

		// Wait for context to be done or timeout
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}

		c.conn.Close()
		c.conn = nil
	}
}
