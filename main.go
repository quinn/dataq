package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.quinn.io/dataq/plugin"
	"go.quinn.io/dataq/queue"
	"go.quinn.io/dataq/worker"
	"gopkg.in/yaml.v3"
)

type Config struct {
	QueuePath string                 `yaml:"queue_path"`
	Plugins   []*plugin.PluginConfig `yaml:"plugins"`
}

func main() {
	// Check if status command
	if len(os.Args) > 1 && os.Args[1] == "status" {
		if err := printStatus(); err != nil {
			log.Printf("Error: %v", err)
			os.Exit(1)
		}
		return
	}

	// Load config
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	log.Printf("Loaded config with %d plugins", len(config.Plugins))
	for _, p := range config.Plugins {
		log.Printf("Plugin: %s (enabled: %v, binary: %s)", p.ID, p.Enabled, p.BinaryPath)
	}

	// Initialize queue
	q, err := queue.NewBoltQueue(queue.WithPath(config.QueuePath))
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	defer q.Close()

	// Create worker
	w := worker.New(q, config.Plugins)

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)
		cancel()
	}()

	// Start worker
	if err := w.Start(ctx); err != nil && err != context.Canceled {
		log.Printf("Worker stopped with error: %v", err)
		os.Exit(1)
	}
}

func printStatus() error {
	// Load config
	configPath := "config.yaml"
	if len(os.Args) > 2 {
		configPath = os.Args[2]
	}

	var config Config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	// Open queue
	q, err := queue.NewBoltQueue(queue.WithPath(config.QueuePath))
	if err != nil {
		return fmt.Errorf("failed to open queue: %v", err)
	}

	// Get tasks by status
	ctx := context.Background()
	statuses := []queue.TaskStatus{
		queue.TaskStatusPending,
		queue.TaskStatusComplete,
		queue.TaskStatusFailed,
	}

	fmt.Println("\nQueue Status:")
	fmt.Println("=============")

	for _, status := range statuses {
		tasks, err := q.List(ctx, status)
		if err != nil {
			return fmt.Errorf("failed to list %s tasks: %v", status, err)
		}

		fmt.Printf("\n%s Tasks: %d\n", status, len(tasks))
		if len(tasks) > 0 {
			fmt.Printf("Latest tasks:\n")
			// Show last 5 tasks
			start := len(tasks) - 5
			if start < 0 {
				start = 0
			}
			for _, task := range tasks[start:] {
				fmt.Printf("  - %s (plugin: %s, created: %s)\n",
					task.ID, task.PluginID, task.CreatedAt.Format(time.RFC3339))
			}
		}
	}

	return nil
}
