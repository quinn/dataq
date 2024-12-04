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

	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/dq"
	"go.quinn.io/dataq/pluginutil"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
)

// Worker handles task processing and plugin execution
type Worker struct {
	queue   queue.Queue
	plugins map[string]*config.Plugin
	dataDir string
	done    chan struct{}
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// New creates a new Worker
func New(q queue.Queue, plugins []*config.Plugin, dataDir string) *Worker {
	pluginMap := make(map[string]*config.Plugin)
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

func (w *Worker) taskToRequest(ctx context.Context, task *queue.Task, messages chan Message) (*pb.PluginRequest, *pb.DataItem, *config.Plugin, error) {

	// Load the data for this task
	data, err := w.loadData(task.Hash)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load data: %w", err)
	}

	var pluginID string
	if data == nil {
		pluginID = strings.Split(task.ID, "_")[0]
	} else {
		pluginID = data.Meta.PluginId
	}

	cfg, ok := w.plugins[pluginID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("plugin not found: %s", pluginID)
	}

	messages <- Message{
		Type: "info",
		Data: "Processing task: " + task.ID,
	}

	if data == nil {
		messages <- Message{
			Type: "info",
			Data: "[input: nil]",
		}
	} else {
		messages <- Message{
			Type: "info",
			Data: "[input: " + data.Meta.Id + "] [hash: " + data.Meta.Hash + "]",
		}

		messages <- Message{
			Type: "protobuf",
			Data: data.Meta.String(),
		}
	}
	// Execute plugin and get response stream
	return &pb.PluginRequest{
		PluginId:  cfg.ID,
		Operation: "extract",
		Config:    cfg.Config,
		Item:      data,
	}, data, cfg, nil
}

// Start begins processing tasks
func (w *Worker) Start(ctx context.Context, messages chan Message) error {
	log.Println("Starting task processing loop")

	tasks := make(chan *queue.Task)
	go func() {
		for {
			task, err := w.queue.Pop(ctx)
			if err != nil {
				messages <- Message{
					Type: "error",
					Data: fmt.Sprintf("Failed to pop task: %v", err),
				}
				return
			}

			if task == nil {
				continue
			}

			tasks <- task
		}
	}()
	w.processRequests(ctx, tasks, messages)
}

// Stop gracefully stops the worker
func (w *Worker) Stop() {
	close(w.done)
}

func (w *Worker) ProcessSingleTask(ctx context.Context, messages chan Message) error {
	result, err := w.processSingleTask(ctx, messages)
	if err != nil {
		return err
	}
	return errors.New(result.Error)
}

// ProcessSingleTask processes a single task and returns the result
func (w *Worker) processRequests(ctx context.Context, tasks chan *queue.Task, messages chan Message) error {
	messages <- Message{
		Type: "info",
		Data: "Processing task",
	}

	for task := range tasks {
		req, data, cfg, err := w.taskToRequest(ctx, task, messages)
		if err != nil {
			task.Status = queue.TaskStatusFailed
			task.Error = err.Error()
			task.UpdatedAt = time.Now()
			w.queue.Update(ctx, task)
			return nil
		}

		responses, err := pluginutil.Execute(ctx, cfg, task, task)
		if err != nil {
			task.Status = queue.TaskStatusFailed
			task.Error = err.Error()
			task.UpdatedAt = time.Now()
			w.queue.Update(ctx, task)
			return nil
		}

		// Process responses
		var lastError string
		for resp := range responses {
			if resp.Error != "" {
				lastError = resp.Error
				messages <- Message{
					Type: "error",
					Data: resp.Error,
				}
				break
			}
			if resp.Item != nil {
				// Store the data item
				hash, err := w.storeData(resp.Item)
				if err != nil {
					messages <- Message{
						Type: "error",
						Data: err.Error(),
					}
					continue
				}

				// Create a new task for this item
				newTask := queue.NewTask(resp.PluginId, resp.Item.Meta.Id, hash)
				if err := w.queue.Push(ctx, newTask); err != nil {
					messages <- Message{
						Type: "error",
						Data: err.Error(),
					}
					continue
				}

				messages <- Message{
					Type: "info",
					Data: "[task: " + newTask.ID + "] [hash: " + resp.Item.Meta.Hash + "] [parent: " + resp.Item.Meta.ParentHash + "]",
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
			return fmt.Errorf("failed to update task: %w", err)
		}

		return nil
	}

	return nil
}

func (q *Worker) storeData(data *pb.DataItem) (string, error) {
	hash := data.Meta.Hash
	f, err := os.Create(filepath.Join(q.dataDir, hash+".dq"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := dq.Write(f, data); err != nil {
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
