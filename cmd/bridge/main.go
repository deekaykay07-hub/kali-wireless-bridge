package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deekaykay07-hub/kali-wireless-bridge/internal/bridge"
)

func main() {
	var (
		connectAddr = flag.String("connect", "", "Remote TUI address (host:port)")
		token       = flag.String("token", "", "Authentication token from the TUI")
		name        = flag.String("name", "", "Friendly name for this bridge (optional)")
	)

	flag.Parse()

	if *connectAddr == "" || *token == "" {
		fmt.Println("Usage: kali-bridge --connect <host:port> --token <token> [--name <name>]")
		os.Exit(1)
	}

	cfg := bridge.Config{
		RemoteAddr: *connectAddr,
		Token:      *token,
		Name:       *name,
	}

	b, err := bridge.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bridge: %v", err)
	}

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nShutting down bridge...")
		b.Close()
		os.Exit(0)
	}()

	fmt.Printf("Connecting to %s as bridge '%s'...\n", *connectAddr, *name)

	if err := b.Run(); err != nil {
		log.Fatalf("Bridge error: %v", err)
	}
}
