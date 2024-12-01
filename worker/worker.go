package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.quinn.io/dataq/dq"
	"go.quinn.io/dataq/pkg/pluginutil"
	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
)

// Worker handles task processing and plugin execution
type Worker struct {
	queue   queue.Queue
	plugins map[string]*plugin.PluginConfig
	dataDir string
	done    chan struct{}
}

// New creates a new Worker
func New(q queue.Queue, plugins []*plugin.PluginConfig, dataDir string) *Worker {
	pluginMap := make(map[string]*plugin.PluginConfig)
	for _, p := range plugins {
		if p.Enabled {
			pluginMap[p.ID] = p
		}
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Failed to create data directory: %v", err)
	}

	return &Worker{
		queue:   q,
		plugins: pluginMap,
		dataDir: dataDir,
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
	return errors.New(result.Error)
}

// ProcessSingleTask processes a single task and returns the result
func (w *Worker) ProcessSingleTask(ctx context.Context) (*queue.Task, error) {
	task, err := w.queue.Pop(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to pop task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("no pending tasks")
	}

	// Load the data for this task
	data, err := w.loadData(task.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	var pluginID string
	if data == nil {
		pluginID = strings.Split(task.ID, "_")[0]
	} else {
		pluginID = data.Meta.PluginId
	}

	// Find the plugin for this task
	plugin, ok := w.plugins[pluginID]
	if !ok {
		task.Status = queue.TaskStatusFailed
		task.Error = fmt.Sprintf("plugin %s not found", pluginID)
		w.queue.Update(ctx, task)
		return task, nil
	}

	// Execute plugin and get response stream
	responses, err := pluginutil.ExecutePlugin(ctx, plugin, data)
	if err != nil {
		task.Status = queue.TaskStatusFailed
		task.Error = err.Error()
		task.UpdatedAt = time.Now()
		w.queue.Update(ctx, task)
		return task, nil
	}

	// Process responses
	var lastError string
	for resp := range responses {
		if resp.Error != "" {
			lastError = resp.Error
			break
		}
		if resp.Item != nil {
			// Store the data item
			hash, err := w.storeData(resp.Item)
			if err != nil {
				// log.Printf("Failed to store data: %v", err)
				continue
			}

			// Create a new task for this item
			newTask := queue.NewTask(resp.PluginId, resp.Item.Meta.Id, hash)
			if err := w.queue.Push(ctx, newTask); err != nil {
				log.Printf("Failed to create task for item: %v", err)
				continue
			}
		}
	}

	// Update task status based on results
	if lastError != "" {
		task.Status = queue.TaskStatusFailed
		task.Error = lastError
	} else {
		task.Status = queue.TaskStatusComplete
	}

	// Update task status
	task.UpdatedAt = time.Now()
	if err := w.queue.Update(ctx, task); err != nil {
		return task, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (q *Worker) storeData(data *pb.DataItem) (string, error) {
	hash := data.Meta.Hash
	f, err := os.Create(filepath.Join(q.dataDir, hash+".dq"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := dq.WriteDataItem(f, data); err != nil {
		return "", err
	}

	return hash, nil
}

func (w *Worker) loadData(hash string) (*pb.DataItem, error) {
	if hash == "" {
		return nil, nil
	}
	filename := filepath.Join(w.dataDir, hash+".dq")

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer f.Close()

	return dq.Read(f)
}
