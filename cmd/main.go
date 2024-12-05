package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/internal/web"
	"go.quinn.io/dataq/queue"
	"go.quinn.io/dataq/tui"
	"go.quinn.io/dataq/worker"
)

func main() {
	boot, err := boot.New()
	if err != nil {
		log.Fatalf("Failed to boot: %v", err)
	}

	queueItems, err := boot.Queue.List(context.Background(), "")
	if err != nil {
		log.Fatalf("Failed to list queue items: %v", err)
	}

	if len(queueItems) == 0 {
		var plugin config.Plugin
		for _, p := range boot.Config.Plugins {
			if p.ID == "gmail" {
				plugin = *p
				break
			}
		}
		initialTask := queue.InitialTask(plugin)
		if err := boot.Queue.Push(context.Background(), initialTask); err != nil {
			log.Printf("Warning: Failed to add initial task: %v", err)
		}
	}

	modeFlag := flag.String("mode", "one", "Queue mode (one|worker|tui|web)")

	switch *modeFlag {
	case "one":
		messages := make(chan worker.Message, 10)
		boot.Worker.ProcessSingleTask(context.Background(), messages)
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
		boot.Worker.Start(context.Background(), messages)
		for msg := range messages {
			fmt.Printf("[MESSAGE] %s: %s\n", msg.Type, msg.Data)
		}
	case "tui":
		// Create and start the TUI
		p := tea.NewProgram(tui.NewModel(boot))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	case "web":
		web.Run()
	}
}
