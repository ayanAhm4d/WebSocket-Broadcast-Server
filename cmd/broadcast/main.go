package main

import (
	"flag"
	"log"
	"os"

	"broadcast/internal/config"
	"broadcast/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: broadcast [start|connect] [options]")
	}

	switch os.Args[1] {
	case "start":
		cmd := flag.NewFlagSet("start", flag.ExitOnError)
		cfg := &config.ServerConfig{}
		cmd.StringVar(&cfg.Port, "port", "8080", "Port to listen on")
		cmd.StringVar(&cfg.Host, "host", "0.0.0.0", "Host to bind to")
		cmd.Parse(os.Args[2:])

		srv := server.NewServer(cfg)
		srv.Start()

	case "connect":
		cmd := flag.NewFlagSet("connect", flag.ExitOnError)
		cfg := &config.ClientConfig{}
		cmd.StringVar(&cfg.ServerAddr, "addr", "localhost:8080", "Server address")
		cmd.StringVar(&cfg.Username, "username", "", "Username for chat")
		cmd.Parse(os.Args[2:])

		client := server.NewClient(cfg)
		client.Connect()

	default:
		log.Fatal("Unknown command. Use 'start' or 'connect'.")
	}
}
