package worker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.quinn.io/dataq/dq"
	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
	"google.golang.org/protobuf/proto"
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

	data, err := w.loadData(task.DataHash)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	plugin, ok := w.plugins[data.Meta.PluginId]
	if !ok {
		task.Status = queue.TaskStatusFailed
		task.Error = fmt.Sprintf("plugin %s not found", data.Meta.PluginId)
		w.queue.Update(ctx, task)
		return task, nil
	}

	// Execute plugin
	pluginResp, err := w.executePlugin(ctx, plugin, data)
	if err != nil {
		task.Status = queue.TaskStatusFailed
		task.Error = err.Error()
	} else {
		task.Status = queue.TaskStatusComplete
	}

	task.UpdatedAt = time.Now()
	if err := w.queue.Update(ctx, task); err != nil {
		return task, fmt.Errorf("failed to update task: %w", err)
	}

	// Only create tasks if we have items in the response
	if task.Status == queue.TaskStatusComplete && len(pluginResp.Items) > 0 {
		// Create a task for each item in the response
		for _, item := range pluginResp.Items {
			hash, err := w.storeData(item)
			if err != nil {
				log.Printf("Failed to store data: %v", err)
				continue
			}
			newTask := queue.NewTask(item.Meta.PluginId, item.Meta.Id, hash)

			if err := w.queue.Push(ctx, newTask); err != nil {
				log.Printf("Failed to create task for item: %v", err)
				continue
			}
		}
	}

	return task, nil
}

func (w *Worker) executePlugin(ctx context.Context, plugin *plugin.PluginConfig, data *pb.DataItem) (*pb.PluginResponse, error) {
	if _, err := os.Stat(plugin.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", plugin.BinaryPath)
	}

	// Create plugin request
	req := &pb.PluginRequest{
		PluginId:  plugin.ID,
		Operation: "extract",
		Config:    plugin.Config,
		Item:      data,
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

func (q *Worker) generateHash(data *pb.DataItem) string {
	// DataItem must exist since this is only called for plugin response items
	h := sha256.New()
	h.Write(data.RawData)
	return hex.EncodeToString(h.Sum(nil))
}

func (q *Worker) storeData(data *pb.DataItem) (string, error) {
	hash := q.generateHash(data)
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
	filename := filepath.Join(w.dataDir, hash+".dq")

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer f.Close()

	return dq.Read(f)
}
