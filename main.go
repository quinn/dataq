package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type PluginConfig struct {
	ID      string            `yaml:"id"`
	Enabled bool              `yaml:"enabled"`
	Config  map[string]string `yaml:"config"`
}

type Config struct {
	Plugins []PluginConfig `yaml:"plugins"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func runPlugin(ctx context.Context, pluginPath string, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	cmd := exec.CommandContext(ctx, pluginPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin: %v", err)
	}

	// Write request to stdin
	reqData, err := protojson.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	if _, err := stdin.Write(reqData); err != nil {
		return nil, fmt.Errorf("failed to write to stdin: %v", err)
	}
	stdin.Close()

	// Read response from stdout
	respData, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdout: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("plugin execution failed: %v", err)
	}

	var resp pb.PluginResponse
	if err := protojson.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("plugin error: %s", resp.Error)
	}

	return &resp, nil
}

func main() {
	// Load configuration
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Process each enabled plugin
	for _, plugin := range config.Plugins {
		if !plugin.Enabled {
			continue
		}

		// Construct plugin binary path
		pluginPath := filepath.Join("cmd", "plugins", plugin.ID, plugin.ID)
		if _, err := os.Stat(pluginPath); err != nil {
			log.Printf("Plugin binary not found: %s", pluginPath)
			continue
		}

		// Configure the plugin
		req := &pb.PluginRequest{
			PluginId:  plugin.ID,
			Config:    plugin.Config,
			Operation: "configure",
		}

		if _, err := runPlugin(ctx, pluginPath, req); err != nil {
			log.Printf("Failed to configure plugin %s: %v", plugin.ID, err)
			continue
		}

		// Extract data
		req.Operation = "extract"
		resp, err := runPlugin(ctx, pluginPath, req)
		if err != nil {
			log.Printf("Failed to extract data from plugin %s: %v", plugin.ID, err)
			continue
		}

		// Process the extracted data
		for _, item := range resp.Items {
			fmt.Printf("Found file: %s\n", item.SourceId)
			fmt.Printf("  Size: %s\n", item.Metadata["size"])
			fmt.Printf("  Last Modified: %s\n", item.Metadata["last_modified"])
			fmt.Println()
		}
	}
}
