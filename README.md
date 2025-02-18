# WebSocket Broadcast Server

A scalable WebSocket broadcast server implementation in Go that enables real-time message broadcasting between multiple clients. This project provides both server and client components with support for multiple concurrent connections, graceful shutdown, and username identification.

## Features

The WebSocket Broadcast Server includes the following key features:

- Real-time message broadcasting to all connected clients
- Support for multiple concurrent client connections
- Username identification for message attribution
- Interactive command-line interface for clients
- Graceful shutdown handling
- Thread-safe connection management
- Configurable server settings
- Clean separation of concerns through modular architecture

## Project Structure

```
broadcast-server/
├── cmd/
│   └── broadcast/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   └── server/
│       ├── server.go
│       └── client.go
├── pkg/
│   └── websocket/
│       └── connection.go
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.16 or higher
- GitHub.com/gorilla/websocket package

## Installation

1. Clone the repository:
```bash
git clone https://github.com/ayanAhm4d/WebSocket-Broadcast-Server.git
cd WebSocket-Broadcast-Server
```


2. Install dependencies:
```bash
go mod tidy
```

## Usage

### Starting the Server

To start the WebSocket server:

```bash
go run cmd/broadcast/main.go start [options]

Options:
  -port string    Port to listen on (default "8080")
  -host string    Host to bind to (default "0.0.0.0")
```

Example:
```bash
go run cmd/broadcast/main.go start -port 8080
```

### Connecting as a Client

To connect to the server as a client:

```bash
go run cmd/broadcast/main.go connect [options]

Options:
  -addr string      Server address (default "localhost:8080")
  -username string  Username for chat (optional)
```

Example:
```bash
go run cmd/broadcast/main.go connect -addr localhost:8080 -username Alice
```

### Client Commands

Once connected as a client:
- Type your message and press Enter to send
- Press Ctrl+C to disconnect
- Messages from other clients will appear automatically
- The prompt '>' indicates you can type a new message

## Configuration

### Server Configuration

The server can be configured using command-line flags or by modifying the `ServerConfig` struct in `internal/config/config.go`:

- Host: The network interface to bind to
- Port: The port number to listen on

### Client Configuration

Client settings can be configured through command-line flags or the `ClientConfig` struct:

- ServerAddr: The address of the WebSocket server
- Username: Optional username for message attribution

## Implementation Details

### Server Component

The server implementation includes:
- Connection management using goroutines
- Thread-safe client tracking
- Message broadcasting to all connected clients
- Graceful shutdown handling
- Error handling and logging

### Client Component

The client implementation provides:
- Interactive command-line interface
- Concurrent message handling
- Connection state management
- Username support
- Clean shutdown handling

## Error Handling

The application implements comprehensive error handling:
- Connection errors are logged and handled appropriately
- Network disconnections are detected and managed
- Resource cleanup is performed during shutdown
- User input errors are handled gracefully

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



## Built With

- [Go](https://golang.org/) - The programming language used
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket implementation for Go



## Acknowledgments

- Gorilla WebSocket team for providing the WebSocket implementation
- Go team for the excellent standard library and tooling

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.
