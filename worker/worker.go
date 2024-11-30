package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"go.quinn.io/dataq/plugin"
	"go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
)

// Worker handles task processing and plugin execution
type Worker struct {
	queue   queue.Queue
	plugins map[string]*plugin.PluginConfig
	done    chan struct{}
}

// New creates a new Worker
func New(q queue.Queue, plugins []*plugin.PluginConfig) *Worker {
	pluginMap := make(map[string]*plugin.PluginConfig)
	for _, p := range plugins {
		if p.Enabled {
			pluginMap[p.ID] = p
		}
	}

	return &Worker{
		queue:   q,
		plugins: pluginMap,
		done:    make(chan struct{}),
	}
}

// Start begins processing tasks
func (w *Worker) Start(ctx context.Context) error {
	log.Println("Starting task processing loop")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.done:
			return nil
		default:
			if err := w.processSingleTask(ctx); err != nil {
				log.Printf("Error processing task: %v", err)
				time.Sleep(time.Second)
			}
		}
	}
}

// Stop gracefully stops the worker
func (w *Worker) Stop() {
	close(w.done)
}

func (w *Worker) processSingleTask(ctx context.Context) error {
	task, err := w.queue.Pop(ctx)
	if err != nil {
		return fmt.Errorf("failed to pop task: %w", err)
	}

	plugin, ok := w.plugins[task.PluginID]
	if !ok {
		task.Status = queue.TaskStatusFailed
		task.Error = fmt.Sprintf("plugin %s not found", task.PluginID)
		return w.queue.Update(ctx, task)
	}

	// Execute plugin
	result, err := w.executePlugin(ctx, plugin, task)
	if err != nil {
		task.Status = queue.TaskStatusFailed
		task.Error = err.Error()
	} else {
		task.Status = queue.TaskStatusComplete
		task.Result = result
	}

	task.UpdatedAt = time.Now()
	if err := w.queue.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Create next task if completed successfully
	if task.Status == queue.TaskStatusComplete {
		nextTask := &queue.Task{
			ID:        fmt.Sprintf("%s_%d", task.PluginID, time.Now().UnixNano()),
			PluginID:  task.PluginID,
			Config:    task.Config,
			Status:    queue.TaskStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := w.queue.Push(ctx, nextTask); err != nil {
			log.Printf("Failed to create next task: %v", err)
		}
	}

	return nil
}

func (w *Worker) executePlugin(ctx context.Context, plugin *plugin.PluginConfig, task *queue.Task) (*proto.PluginResponse, error) {
	if _, err := os.Stat(plugin.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", plugin.BinaryPath)
	}

	cmd := exec.CommandContext(ctx, plugin.BinaryPath)
	cmd.Env = os.Environ()
	for k, v := range task.Config {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w: %s", err, string(output))
	}

	// For now, just create an empty response
	// In a real implementation, we would parse the plugin output
	return &proto.PluginResponse{
		PluginId: plugin.ID,
	}, nil
}
