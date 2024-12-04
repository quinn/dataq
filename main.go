package main

import (
	"context"
	"encoding/json"
	"flag"
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

	queueItems, err := q.List(context.Background(), "")
	if err != nil {
		log.Fatalf("Failed to list queue items: %v", err)
	}

	if len(queueItems) == 0 {
		var plugin config.Plugin
		for _, p := range cfg.Plugins {
			if p.ID == "gmail" {
				plugin = *p
				break
			}
		}
		initialTask := queue.InitialTask(plugin)
		if err := q.Push(context.Background(), initialTask); err != nil {
			log.Printf("Warning: Failed to add initial task: %v", err)
		}
	}

	// Create worker
	w := worker.New(q, cfg.Plugins, cfg.DataDir)

	modeFlag := flag.String("mode", "one", "Queue mode (one|worker|tui)")

	switch *modeFlag {
	case "one":
		messages := make(chan worker.Message, 10)
		w.ProcessSingleTask(context.Background(), messages)
		for msg := range messages {
			// fmt.Printf("[MESSAGE] %s: %s\n", msg.Type, msg.Data)

			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}
			fmt.Printf("%s\n", jsonData)
		}
	case "worker":
		messages := make(chan worker.Message, 10)
		w.Start(context.Background(), messages)
		for msg := range messages {
			fmt.Printf("[MESSAGE] %s: %s\n", msg.Type, msg.Data)
		}
	case "tui":

		// Create and start the TUI
		p := tea.NewProgram(tui.NewModel(q, w))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	}
}
