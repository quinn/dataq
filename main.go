package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.quinn.io/dataq/plugin"
	"go.quinn.io/dataq/plugin/filescan"
)

func main() {
	// Create a new plugin registry
	registry := plugin.NewRegistry()

	// Create and register the filesystem scanner plugin
	fsPlugin := filescan.New()
	if err := registry.Register(fsPlugin); err != nil {
		log.Fatal(err)
	}

	// Configure the plugin
	err := fsPlugin.Configure(map[string]string{
		"root_path": ".",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Extract data using the plugin
	items, err := fsPlugin.Extract(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Process the extracted data
	for item := range items {
		fmt.Printf("Found file: %s\n", item.SourceId)
		fmt.Printf("  Size: %s\n", item.Metadata["size"])
		fmt.Printf("  Last Modified: %s\n", item.Metadata["last_modified"])
		fmt.Println()
	}
}
