package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jaysinh001/Glide-Go/internal/discovery"
	"github.com/Jaysinh001/Glide-Go/internal/engine"
	"github.com/Jaysinh001/Glide-Go/internal/input"
)

func main() {
	log.Println("Starting Glide Engine (Phase 0)")

	// Root context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start mDNS broadcaster (Phase 1.1)
	_, err := discovery.StartMDNS(ctx, 50506, 50505)
	if err != nil {
		log.Println("mDNS failed to start:", err)
	}

	// Handle Ctrl+C / SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutdown signal received")
		cancel()
	}()

	// Hardcoded UDP address for Phase 0
	listener, err := engine.NewUDPListener(":50505")
	if err != nil {
		log.Fatal("Failed to start UDP listener:", err)
	}

	// Create Windows native injector
	injector := input.NewInjector()

	// Create and run engine
	glideEngine := engine.NewEngine(listener, injector)
	if err := glideEngine.Run(ctx); err != nil {
		log.Fatal("Engine stopped with error:", err)
	}

	log.Println("Glide Engine stopped cleanly")
}
