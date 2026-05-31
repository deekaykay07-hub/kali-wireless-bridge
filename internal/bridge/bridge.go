package bridge

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
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

// Message types (simple protocol)
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type CommandPayload struct {
	ID         string   `json:"id"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	Privileged bool     `json:"privileged"`
}

type OutputPayload struct {
	ID     string `json:"id"`
	Stream string `json:"stream"` // "stdout" or "stderr"
	Data   string `json:"data"`
}

type ResultPayload struct {
	ID       string `json:"id"`
	ExitCode int    `json:"exit_code"`
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

	var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Bad message: %s", message)
			continue
		}

		if msg.Type == "command" {
			var cmd CommandPayload
			json.Unmarshal(msg.Payload, &cmd)
			go b.handleCommand(cmd)
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

// handleCommand runs a command and streams output back
func (b *Bridge) handleCommand(cmd CommandPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fullCmd := append([]string{cmd.Command}, cmd.Args...)
	log.Printf("Executing: %v (privileged=%v)", fullCmd, cmd.Privileged)

	var execCmd *exec.Cmd
	if cmd.Privileged {
		execCmd = exec.CommandContext(ctx, "sudo", fullCmd...)
	} else {
		execCmd = exec.CommandContext(ctx, cmd.Command, cmd.Args...)
	}

	stdoutPipe, _ := execCmd.StdoutPipe()
	stderrPipe, _ := execCmd.StderrPipe()

	if err := execCmd.Start(); err != nil {
		b.sendError(cmd.ID, err.Error())
		return
	}

	// Stream stdout
	go b.streamOutput(cmd.ID, "stdout", stdoutPipe)
	// Stream stderr
	go b.streamOutput(cmd.ID, "stderr", stderrPipe)

	err := execCmd.Wait()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	// Send final result
	result := ResultPayload{
		ID:       cmd.ID,
		ExitCode: exitCode,
	}
	b.conn.WriteJSON(map[string]interface{}{
		"type":    "result",
		"payload": result,
	})
}

func (b *Bridge) streamOutput(id string, stream string, pipe io.ReadCloser) {
	defer pipe.Close()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		payload := OutputPayload{
			ID:     id,
			Stream: stream,
			Data:   scanner.Text() + "\n",
		}
		b.conn.WriteJSON(map[string]interface{}{
			"type":    "output",
			"payload": payload,
		})
	}
}

func (b *Bridge) sendError(id string, msg string) {
	b.conn.WriteJSON(map[string]interface{}{
		"type": "error",
		"payload": map[string]string{"id": id, "message": msg},
	})
}

// Close shuts down the bridge
func (b *Bridge) Close() error {
	close(b.done)
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}
