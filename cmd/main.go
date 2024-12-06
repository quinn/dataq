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
	"go.quinn.io/dataq/internal/web"
	"go.quinn.io/dataq/tui"
	"go.quinn.io/dataq/worker"
)

func main() {
	b, err := boot.New()
	if err != nil {
		log.Fatalf("Failed to boot: %v", err)
	}

	modeFlag := flag.String("mode", "", "Queue mode (one|worker|tui|web)")
	flag.Parse()

	switch *modeFlag {
	case "one":
		messages := make(chan worker.Message, 10)
		if _, err := b.Worker.ProcessSingleTask(context.Background(), messages); err != nil {
			log.Printf("Failed to process single task: %v", err)
			break
		}
		for msg := range messages {
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}
			fmt.Printf("%s\n", jsonData)
		}
	case "worker":
		messages := make(chan worker.Message, 10)
		b.Worker.Start(context.Background(), messages)
		for msg := range messages {
			fmt.Printf("[MESSAGE] %s: %s\n", msg.Type, msg.Data)
		}
	case "tui":
		// Create and start the TUI
		p := tea.NewProgram(tui.NewModel(b))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	case "web":
		web.Run(b)
	default:
		log.Fatalf("Unknown mode: \"%s\"", *modeFlag)
	}
}
