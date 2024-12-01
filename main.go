package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.quinn.io/dataq/plugin"
	"go.quinn.io/dataq/queue"
	"go.quinn.io/dataq/tui"
	"go.quinn.io/dataq/worker"
	"gopkg.in/yaml.v3"
)

type Config struct {
	QueuePath string                 `yaml:"queue_path"`
	Plugins   []*plugin.PluginConfig `yaml:"plugins"`
}

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

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Initialize queue
	q, err := queue.NewFileQueue(config.QueuePath)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	defer q.Close()

	// Add initial task for Gmail plugin
	initialTask := &queue.Task{
		ID:        "gmail_initial",
		PluginID:  "gmail",
		Status:    queue.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config:    config.Plugins[1].Config,  // Use config from the Gmail plugin
		Data:      nil,  // Initial task has no data
	}
	if err := q.Push(context.Background(), initialTask); err != nil {
		log.Printf("Warning: Failed to add initial task: %v", err)
	}

	// Create worker
	w := worker.New(q, config.Plugins)

	// Create and start the TUI
	p := tea.NewProgram(tui.NewModel(q, w))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
