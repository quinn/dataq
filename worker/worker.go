package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
	"google.golang.org/protobuf/proto"
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
			log.Println("Context cancelled, stopping task processing")
			return ctx.Err()
		case <-w.done:
			log.Println("Worker stopped, exiting task processing")
			return nil
		default:
			if err := w.processSingleTask(ctx); err != nil {
				log.Printf("Error processing task: %v", err)
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

	plugin, ok := w.plugins[task.Data.Meta.PluginId]
	if !ok {
		task.Meta.Status = queue.TaskStatusFailed
		task.Meta.Error = fmt.Sprintf("plugin %s not found", task.Data.Meta.PluginId)
		result.Error = fmt.Errorf("plugin not found: %s", task.Data.Meta.PluginId)
		w.queue.Update(ctx, &task.Meta)
		return result, nil
	}

	// Execute plugin
	pluginResp, err := w.executePlugin(ctx, plugin, task)
	if err != nil {
		task.Meta.Status = queue.TaskStatusFailed
		task.Meta.Error = err.Error()
		result.Error = err
	} else {
		task.Meta.Status = queue.TaskStatusComplete
	}

	task.Meta.UpdatedAt = time.Now()
	if err := w.queue.Update(ctx, &task.Meta); err != nil {
		result.Error = fmt.Errorf("failed to update task: %w", err)
		return result, nil
	}

	// Only create tasks if we have items in the response
	if task.Meta.Status == queue.TaskStatusComplete && len(pluginResp.Items) > 0 {
		// Create a task for each item in the response
		for _, item := range pluginResp.Items {
			newTask := queue.NewTask(item)

			if err := w.queue.Push(ctx, newTask); err != nil {
				log.Printf("Failed to create task for item: %v", err)
				continue
			}
		}
	}

	return result, nil
}

func (w *Worker) executePlugin(ctx context.Context, plugin *plugin.PluginConfig, task *queue.Task) (*pb.PluginResponse, error) {
	if _, err := os.Stat(plugin.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", plugin.BinaryPath)
	}

	// Create plugin request
	req := &pb.PluginRequest{
		PluginId:  plugin.ID,
		Operation: "extract",
		Config:    plugin.Config,
		Item:      task.Data,
	}

	// Serialize request using protobuf
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, plugin.BinaryPath)
	cmd.Stdin = bytes.NewReader(reqData)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run plugin
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w\nstderr: %s", err, stderr.String())
	}

	// Parse response using protobuf
	var resp pb.PluginResponse
	if err := proto.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}
