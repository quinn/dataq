package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.quinn.io/dataq/plugin"
	"go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
	"gopkg.in/yaml.v3"
)

type Config struct {
	QueuePath string                 `yaml:"queue_path"`
	Plugins   []*plugin.PluginConfig `yaml:"plugins"`
}

func main() {
	// Check if status command
	if len(os.Args) > 1 && os.Args[1] == "status" {
		if err := printStatus(); err != nil {
			log.Printf("Error: %v", err)
			os.Exit(1)
		}
		return
	}

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

	log.Printf("Loaded config with %d plugins", len(config.Plugins))
	for _, p := range config.Plugins {
		log.Printf("Plugin: %s (enabled: %v, binary: %s)", p.ID, p.Enabled, p.BinaryPath)
	}

	// Initialize queue
	q, err := queue.NewBoltQueue(queue.WithPath(config.QueuePath))
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	defer q.Close()

	// Process each plugin
	for _, p := range config.Plugins {
		if !p.Enabled {
			log.Printf("Plugin %s is disabled", p.ID)
			continue
		}

		// Check if plugin binary exists
		if _, err := os.Stat(p.BinaryPath); err != nil {
			log.Printf("Plugin binary not found: %s", p.BinaryPath)
			continue
		}

		log.Printf("Creating initial task for plugin %s", p.ID)

		// Create initial task
		task := &queue.Task{
			ID:        fmt.Sprintf("%s_%d", p.ID, time.Now().UnixNano()),
			PluginID:  p.ID,
			Config:    p.Config,
			Status:    queue.TaskStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := q.Push(context.Background(), task); err != nil {
			log.Printf("Failed to push task for plugin %s: %v", p.ID, err)
			continue
		}

		log.Printf("Initial task created for plugin %s with ID %s", p.ID, task.ID)
	}

	log.Println("Starting task processing loop")

	// Process tasks
	for {
		task, err := q.Pop(context.Background())
		if err != nil {
			log.Printf("Failed to pop task: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if task == nil {
			log.Println("No pending tasks")
			time.Sleep(time.Second)
			continue
		}

		// Find plugin config
		var pluginConfig *plugin.PluginConfig
		for _, p := range config.Plugins {
			if p.ID == task.PluginID {
				pluginConfig = p
				break
			}
		}

		if pluginConfig == nil {
			log.Printf("Plugin config not found for task %s", task.ID)
			task.Status = queue.TaskStatusFailed
			task.Error = "plugin config not found"
			if err := q.Update(context.Background(), task); err != nil {
				log.Printf("Failed to update task %s: %v", task.ID, err)
			}
			continue
		}

		// Create plugin request
		pluginReq := &proto.PluginRequest{
			PluginId: task.PluginID,
			Config:   task.Config,
		}
		if task.Data != nil {
			pluginReq.Item = task.Data
		}

		// Run plugin
		plugin, err := plugin.NewPlugin(pluginConfig)
		if err != nil {
			log.Printf("Failed to create plugin %s: %v", pluginConfig.ID, err)
			task.Status = queue.TaskStatusFailed
			task.Error = fmt.Sprintf("failed to create plugin: %v", err)
			if err := q.Update(context.Background(), task); err != nil {
				log.Printf("Failed to update task %s: %v", task.ID, err)
			}
			continue
		}

		resp, err := plugin.Extract(context.Background(), pluginReq)
		if err != nil {
			log.Printf("Plugin %s failed: %v", pluginConfig.ID, err)
			task.Status = queue.TaskStatusFailed
			task.Error = fmt.Sprintf("plugin failed: %v", err)
			if err := q.Update(context.Background(), task); err != nil {
				log.Printf("Failed to update task %s: %v", task.ID, err)
			}
			continue
		}

		// Check if plugin returned any data
		if len(resp.Items) == 0 {
			log.Printf("Plugin %s returned no data", pluginConfig.ID)
			task.Status = queue.TaskStatusFailed
			task.Error = "plugin returned no data"
			if err := q.Update(context.Background(), task); err != nil {
				log.Printf("Failed to update task %s: %v", task.ID, err)
			}
			continue
		}

		// Update task with result
		task.Status = queue.TaskStatusComplete
		task.Result = resp
		task.UpdatedAt = time.Now()

		if err := q.Update(context.Background(), task); err != nil {
			log.Printf("Failed to update task %s: %v", task.ID, err)
			continue
		}

		// Create next task for each item
		for _, item := range resp.Items {
			nextTask := &queue.Task{
				ID:        fmt.Sprintf("%s_%d", pluginConfig.ID, time.Now().UnixNano()),
				PluginID:  pluginConfig.ID,
				Config:    task.Config,
				Data:      item,
				Status:    queue.TaskStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := q.Push(context.Background(), nextTask); err != nil {
				log.Printf("Failed to push next task for plugin %s: %v", pluginConfig.ID, err)
			}
		}
	}
}

func printStatus() error {
	// Load config
	configPath := "config.yaml"
	if len(os.Args) > 2 {
		configPath = os.Args[2]
	}

	var config Config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	// Open queue
	q, err := queue.NewBoltQueue(queue.WithPath(config.QueuePath))
	if err != nil {
		return fmt.Errorf("failed to open queue: %v", err)
	}

	// Get tasks by status
	ctx := context.Background()
	statuses := []queue.TaskStatus{
		queue.TaskStatusPending,
		queue.TaskStatusComplete,
		queue.TaskStatusFailed,
	}

	fmt.Println("\nQueue Status:")
	fmt.Println("=============")

	for _, status := range statuses {
		tasks, err := q.List(ctx, status)
		if err != nil {
			return fmt.Errorf("failed to list %s tasks: %v", status, err)
		}

		fmt.Printf("\n%s Tasks: %d\n", status, len(tasks))
		if len(tasks) > 0 {
			fmt.Printf("Latest tasks:\n")
			// Show last 5 tasks
			start := len(tasks) - 5
			if start < 0 {
				start = 0
			}
			for _, task := range tasks[start:] {
				fmt.Printf("  - %s (plugin: %s, created: %s)\n",
					task.ID, task.PluginID, task.CreatedAt.Format(time.RFC3339))
			}
		}
	}

	return nil
}
