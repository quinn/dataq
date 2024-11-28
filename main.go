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

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-boltdb/v2"
	"github.com/ThreeDotsLabs/watermill/message"
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
	QueuePath string `yaml:"queue_path"`
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

	// Initialize queue
	if config.QueuePath == "" {
		config.QueuePath = "queue.db"
	}
	q, err := watermill.NewBoltDBQueue(watermill.NewBoltDBConfig(config.QueuePath))
	if err != nil {
		log.Fatalf("Failed to initialize queue: %v", err)
	}
	defer q.Close()

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

		// Create initial task
		task := &message.Message{
			UUID:        fmt.Sprintf("%s_%d", plugin.ID, time.Now().UnixNano()),
			Payload:     []byte(fmt.Sprintf("%+v", plugin.Config)),
			Metadata:    map[string]string{"plugin_id": plugin.ID},
			ContentType: "application/json",
		}

		if err := q.Publish(ctx, "tasks", task); err != nil {
			log.Printf("Failed to queue task for plugin %s: %v", plugin.ID, err)
			continue
		}
	}

	// Process tasks
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := q.Get(ctx, "tasks")
		if err != nil {
			log.Printf("Failed to get task: %v", err)
			continue
		}
		if msg == nil {
			// No more tasks
			return
		}

		// Configure the plugin
		req := &pb.PluginRequest{
			PluginId:  msg.Metadata["plugin_id"],
			Config:    map[string]string(msg.Payload),
			Operation: "configure",
		}

		if _, err := runPlugin(ctx, req.PluginId, req); err != nil {
			log.Printf("Failed to configure plugin %s: %v", req.PluginId, err)
			q.Ack(ctx, msg.UUID)
			continue
		}

		// Extract data
		req.Operation = "extract"
		resp, err := runPlugin(ctx, req.PluginId, req)
		if err != nil {
			log.Printf("Failed to extract data from plugin %s: %v", req.PluginId, err)
			q.Ack(ctx, msg.UUID)
			continue
		}

		// Process the extracted data
		for _, item := range resp.Items {
			fmt.Printf("Found item from %s: %s\n", item.PluginId, item.SourceId)
			for k, v := range item.Metadata {
				fmt.Printf("  %s: %s\n", k, v)
			}
			fmt.Println()

			// If plugin returned a next page token, queue a new task for it
			if nextToken, ok := item.Metadata["next_page_token"]; ok && nextToken != "" {
				nextConfig := make(map[string]string)
				for k, v := range req.Config {
					nextConfig[k] = v
				}
				nextConfig["page_token"] = nextToken

				nextTask := &message.Message{
					UUID:        fmt.Sprintf("%s_%d", req.PluginId, time.Now().UnixNano()),
					Payload:     []byte(fmt.Sprintf("%+v", nextConfig)),
					Metadata:    map[string]string{"plugin_id": req.PluginId},
					ContentType: "application/json",
				}

				if err := q.Publish(ctx, "tasks", nextTask); err != nil {
					log.Printf("Failed to queue next page task: %v", err)
				}
			}
		}

		// Acknowledge task
		if err := q.Ack(ctx, msg.UUID); err != nil {
			log.Printf("Failed to acknowledge task: %v", err)
		}
	}
}
