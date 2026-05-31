package bridge

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// Config holds bridge configuration
type Config struct {
	RemoteAddr string
	Token      string
	Name       string
}

// Bridge represents the local hardware bridge agent
type Bridge struct {
	config Config
	conn   *websocket.Conn
	done   chan struct{}
}

// New creates a new Bridge instance
func New(cfg Config) (*Bridge, error) {
	if cfg.Name == "" {
		hostname, _ := os.Hostname()
		cfg.Name = hostname
	}

	return &Bridge{
		config: cfg,
		done:   make(chan struct{}),
	}, nil
}

// Run connects to the remote TUI and starts the bridge loop
func (b *Bridge) Run() error {
	u := url.URL{Scheme: "ws", Host: b.config.RemoteAddr, Path: "/bridge"}

	header := map[string][]string{
		"Authorization": {"Bearer " + b.config.Token},
		"X-Bridge-Name": {b.config.Name},
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return fmt.Errorf("failed to connect to remote TUI: %w", err)
	}
	b.conn = conn

	log.Printf("Connected to remote TUI as bridge '%s'", b.config.Name)

	// Send initial capabilities
	if err := b.sendCapabilities(); err != nil {
		return err
	}

	// Main message loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		// TODO: Handle incoming commands from TUI
		log.Printf("Received message: %s", message)

		// Example: echo back for now
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}

	return nil
}

func (b *Bridge) sendCapabilities() error {
	caps := map[string]interface{}{
		"type":     "capabilities",
		"name":     b.config.Name,
		"version": "0.1.0",
		"features": []string{
			"execute_command",
			"list_wireless_interfaces",
			"stream_output",
		},
	}

	return b.conn.WriteJSON(caps)
}

// Close shuts down the bridge
func (b *Bridge) Close() error {
	close(b.done)
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}

// TODO: Implement real hardware detection and command execution
// This will live in internal/hardware/
