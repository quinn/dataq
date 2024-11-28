package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v3"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type PluginConfig struct {
	ID         string            `yaml:"id"`
	Enabled    bool              `yaml:"enabled"`
	BinaryPath string            `yaml:"binary_path"` // Path to the plugin binary
	Config     map[string]string `yaml:"config"`
}

type Config struct {
	Plugins []struct {
		ID         string            `yaml:"id"`
		Enabled    bool              `yaml:"enabled"`
		BinaryPath string            `yaml:"binary_path"`
		Config     map[string]string `yaml:"config"`
	} `yaml:"plugins"`
	StateFile string `yaml:"state_file"`
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
	log.Printf("Running plugin %s with config: %+v", pluginPath, req.Config)
	cmd := exec.CommandContext(ctx, pluginPath)
	cmd.Stderr = os.Stderr

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

	// Initialize state manager
	if config.StateFile == "" {
		config.StateFile = "plugin_state.json"
	}
	stateManager, err := plugin.NewStateManager(config.StateFile)
	if err != nil {
		log.Fatalf("Failed to initialize state manager: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Process each enabled plugin
	for _, plugin := range config.Plugins {
		if !plugin.Enabled {
			continue
		}

		// Check if plugin binary exists
		if _, err := os.Stat(plugin.BinaryPath); err != nil {
			log.Printf("Plugin binary not found: %s", plugin.BinaryPath)
			continue
		}

		// Get plugin state
		state := stateManager.GetState(plugin.ID)

		// Merge state metadata into plugin config
		pluginConfig := make(map[string]string)
		for k, v := range plugin.Config {
			pluginConfig[k] = v
		}
		for k, v := range state.Metadata {
			pluginConfig[k] = v
		}

		log.Printf("Processing plugin %s with config: %+v", plugin.ID, pluginConfig)

		// Configure the plugin
		req := &pb.PluginRequest{
			PluginId:  plugin.ID,
			Config:    pluginConfig,
			Operation: "configure",
		}

		if _, err := runPlugin(ctx, plugin.BinaryPath, req); err != nil {
			log.Printf("Failed to configure plugin %s: %v", plugin.ID, err)
			continue
		}

		// Extract data
		req.Operation = "extract"
		resp, err := runPlugin(ctx, plugin.BinaryPath, req)
		if err != nil {
			log.Printf("Failed to extract data from plugin %s: %v", plugin.ID, err)
			continue
		}

		// Process the extracted data
		for _, item := range resp.Items {
			fmt.Printf("Found item from %s: %s\n", item.PluginId, item.SourceId)
			for k, v := range item.Metadata {
				fmt.Printf("  %s: %s\n", k, v)
			}
			fmt.Println()

			// Update plugin state with metadata from the last item
			if nextToken, ok := item.Metadata["next_page_token"]; ok {
				state.Metadata["page_token"] = nextToken
			}
		}

		// Update plugin state
		state.LastRun = time.Now().Unix()
		if err := stateManager.UpdateState(state); err != nil {
			log.Printf("Failed to update plugin state: %v", err)
		}
	}
}
