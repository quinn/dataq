package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/internal/web"
)

func main() {
	// Initialize the boot package
	b, err := boot.New()
	if err != nil {
		log.Fatalf("Failed to initialize boot: %v", err)
	}

	// Create a root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start all plugins
	if err := b.StartPlugins(); err != nil {
		log.Fatalf("Failed to start plugins: %v", err)
	}

	// Create a WaitGroup for coordinating shutdown
	var wg sync.WaitGroup

	// Start the web server in a goroutine
	wg.Add(1)
	webServer := web.NewServer(b)
	go func() {
		defer wg.Done()
		if err := webServer.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Printf("Web server error: %v", err)
		}
	}()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")

	// Cancel context to notify all components
	cancel()

	// Initiate graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown web server first
	if err := webServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Web server shutdown error: %v", err)
	}

	// Then shutdown boot components (including plugins)
	if err := b.Shutdown(shutdownCtx); err != nil {
		log.Printf("Boot shutdown error: %v", err)
	}

	// Wait for all components to shut down
	wg.Wait()
	log.Println("Shutdown complete")
}
