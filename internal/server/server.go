package server

import (
	"broadcast/internal/config"
	"broadcast/pkg/websocket"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Server represents the websocket server
type Server struct {
	cfg     *config.ServerConfig
	clients map[*websocket.Connection]bool
	mutex   sync.RWMutex

	//Channels for sever operations
	broadcast  chan []byte
	register   chan *websocket.Connection
	unregister chan *websocket.Connection
}

// NewServer initializes a new websocket server
func NewServer(cfg *config.ServerConfig) *Server {
	return &Server{
		cfg:        cfg,
		clients:    make(map[*websocket.Connection]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *websocket.Connection),
		unregister: make(chan *websocket.Connection),
	}
}

// Handling Incoming Messages
func (s *Server) handleMessages(client *websocket.Connection) {
	defer func() {
		s.unregister <- client
	}()

	for {
		msg, err := client.ReadMEssage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		s.broadcast <- msg
	}
}

// Handling websocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	client := websocket.NewConnection(conn)
	s.register <- client
	go s.handleMessages(client)
}

// Starts the websocket server
func (s *Server) Start() {
	http.HandleFunc("/ws", s.handleWebSocket)
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	//Handle shutdown gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go s.run()
	go func() {
		log.Printf("Server started on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.shutdown(ctx)

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// run handles the websocket events in loop
func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				client.Close()
			}
			s.mutex.Unlock()
		case message := <-s.broadcast:
			s.broadcastMessage(message)
		}
	}
}

func (s *Server) broadcastMessage(message []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for client := range s.clients {
		go func(c *websocket.Connection) {
			if err := c.WriteMessage(message); err != nil {
				log.Printf("Write error: %v", err)
				s.unregister <- c
			}
		}(client)
	}
}

func (s *Server) shutdown(ctx context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for client := range s.clients {
		client.Close()
	}
}
