package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Connection wrap a websocket connection with synchronization
type Connection struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

// NewConnection creates a new wrapped websocket connection
func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

// Upgrade upgrading and HTTP Connection to websocket
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

// Reading Messages
func (c *Connection) ReadMEssage() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	return message, err
}

// Writing Message
func (c *Connection) WriteMessage(message []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, message)
}

// Closing the websocket connection
func (c *Connection) Close() error {
	return c.conn.Close()
}
