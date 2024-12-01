package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"go.quinn.io/dataq/plugin"
	"go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
	"google.golang.org/protobuf/encoding/protojson"
)

// TaskResult represents the result of processing a single task
type TaskResult struct {
	Task   *queue.Task
	Error  error
	Output string
}

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
	result, err := w.ProcessSingleTask(ctx)
	if err != nil {
		return err
	}
	return result.Error
}

// ProcessSingleTask processes a single task and returns the result
func (w *Worker) ProcessSingleTask(ctx context.Context) (*TaskResult, error) {
	task, err := w.queue.Pop(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to pop task: %w", err)
	}

	if task == nil {
		return nil, fmt.Errorf("no pending tasks")
	}

	result := &TaskResult{
		Task: task,
	}

	plugin, ok := w.plugins[task.PluginID]
	if !ok {
		task.Status = queue.TaskStatusFailed
		task.Error = fmt.Sprintf("plugin %s not found", task.PluginID)
		result.Error = fmt.Errorf("plugin not found: %s", task.PluginID)
		w.queue.Update(ctx, task)
		return result, nil
	}

	// Execute plugin
	pluginResp, err := w.executePlugin(ctx, plugin, task)
	if err != nil {
		task.Status = queue.TaskStatusFailed
		task.Error = err.Error()
		result.Error = err
	} else {
		task.Status = queue.TaskStatusComplete
		task.Result = pluginResp
	}

	task.UpdatedAt = time.Now()
	if err := w.queue.Update(ctx, task); err != nil {
		result.Error = fmt.Errorf("failed to update task: %w", err)
		return result, nil
	}

	// Only create tasks if we have items in the response
	if task.Status == queue.TaskStatusComplete && len(pluginResp.Items) > 0 {
		// Create a task for each item in the response
		for _, item := range pluginResp.Items {
			newTask := &queue.Task{
				ID:        fmt.Sprintf("%s_%d", task.PluginID, time.Now().UnixNano()),
				PluginID:  task.PluginID,
				Status:    queue.TaskStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Config:    task.Config,
				Data:      item,
			}

			if err := w.queue.Push(ctx, newTask); err != nil {
				log.Printf("Failed to create task for item: %v", err)
				continue
			}
		}
	}

	return result, nil
}

func (w *Worker) executePlugin(ctx context.Context, plugin *plugin.PluginConfig, task *queue.Task) (*proto.PluginResponse, error) {
	if _, err := os.Stat(plugin.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", plugin.BinaryPath)
	}

	// Create plugin request
	req := &proto.PluginRequest{
		PluginId:  plugin.ID,
		Operation: "extract",
		Config:    plugin.Config,
		Item:      task.Data,
	}

	// Marshal request to JSON using protojson
	reqData, err := protojson.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create command with stdin pipe
	cmd := exec.CommandContext(ctx, plugin.BinaryPath)
	cmd.Env = os.Environ()
	for k, v := range plugin.Config {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Get stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	// Create buffer for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin: %w", err)
	}

	// Write request to stdin
	if _, err := io.Copy(stdin, bytes.NewReader(reqData)); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}
	stdin.Close()

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w: %s", err, stderr.String())
	}

	// Parse response using protojson
	var resp proto.PluginResponse
	if err := protojson.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}
