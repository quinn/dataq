package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/internal/web"
)

func main() {
	// Initialize the boot package
	b, err := boot.New()
	if err != nil {
		log.Fatalf("Failed to initialize boot: %v", err)
	}

	// Start all plugins
	if err := b.StartPlugins(); err != nil {
		log.Fatalf("Failed to start plugins: %v", err)
	}

	// Create a WaitGroup for coordinating shutdown
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the web server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		web.Run(b)
	}()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
	
	// Cancel context and shutdown all components
	cancel()
	b.Shutdown()

	// Wait for all components to shut down
	wg.Wait()
}
