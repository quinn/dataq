package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/queue"
	"go.quinn.io/dataq/tui"
	"go.quinn.io/dataq/worker"
	"gopkg.in/yaml.v3"
)

func main() {
	// Load config
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Initialize queue
	q, err := queue.NewQueue("file", cfg.DataDir)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	defer q.Close()

	initialTask := queue.InitialTask("gmail")
	if err := q.Push(context.Background(), initialTask); err != nil {
		log.Printf("Warning: Failed to add initial task: %v", err)
	}

	// Create worker
	w := worker.New(q, cfg.Plugins, cfg.DataDir)

	// Create and start the TUI
	p := tea.NewProgram(tui.NewModel(q, w))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
