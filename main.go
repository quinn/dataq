package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v3"

	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
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

func runPlugin(_ context.Context, binaryPath string, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	log.Printf("Running plugin %s with config: %+v", req.PluginId, req.Config)

	// Serialize request to JSON
	reqJson, err := protojson.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Start plugin process
	cmd := exec.Command(binaryPath)
	cmd.Stdin = bytes.NewReader(reqJson)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to start plugin: %v\nstderr: %s", err, stderr.String())
	}

	// Parse response
	resp := &pb.PluginResponse{}
	if err := protojson.Unmarshal(stdout.Bytes(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return resp, nil
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
	q, err := queue.NewBoltQueue(queue.WithPath(config.QueuePath))
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
		task := &queue.Task{
			ID:        fmt.Sprintf("%s_%d", plugin.ID, time.Now().UnixNano()),
			PluginID:  plugin.ID,
			Config:    plugin.Config,
			Status:    queue.TaskStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// Add binary path to config
		task.Config["binary_path"] = plugin.BinaryPath

		if err := q.Push(ctx, task); err != nil {
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

		task, err := q.Pop(ctx)
		if err != nil {
			log.Printf("Failed to pop task: %v", err)
			continue
		}
		if task == nil {
			// No more tasks
			return
		}

		// Configure the plugin
		req := &pb.PluginRequest{
			PluginId:  task.PluginID,
			Config:    task.Config,
			Operation: "configure",
		}

		if _, err := runPlugin(ctx, task.Config["binary_path"], req); err != nil {
			log.Printf("Failed to configure plugin %s: %v", task.PluginID, err)
			q.Update(ctx, task.ID, queue.TaskStatusFailed, nil, err.Error())
			continue
		}

		// Extract data
		req.Operation = "extract"
		resp, err := runPlugin(ctx, task.Config["binary_path"], req)
		if err != nil {
			log.Printf("Failed to extract data from plugin %s: %v", task.PluginID, err)
			q.Update(ctx, task.ID, queue.TaskStatusFailed, nil, err.Error())
			continue
		}

		// Process the extracted data
		for _, item := range resp.Items {
			fmt.Printf("Found item from %s: %s\n", item.PluginId, item.SourceId)

			// Create next task with the response data
			nextConfig := make(map[string]string)
			for k, v := range task.Config {
				nextConfig[k] = v
			}
			// Add all metadata to next task's config
			for k, v := range item.Metadata {
				nextConfig[k] = v
			}

			nextTask := &queue.Task{
				ID:        fmt.Sprintf("%s_%d", task.PluginID, time.Now().UnixNano()),
				PluginID:  task.PluginID,
				Status:    queue.TaskStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := q.Push(ctx, nextTask); err != nil {
				log.Printf("Failed to queue next task: %v", err)
			}
		}

		// Update task status
		if err := q.Update(ctx, task.ID, queue.TaskStatusCompleted, resp, ""); err != nil {
			log.Printf("Failed to update task status: %v", err)
		}
	}
}
