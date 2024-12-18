package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	Type   string `json:"type"`
	Data   string `json:"data"`
	Closed bool   `json:"closed"`
	Done   bool   `json:"done"`
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

	return w.processRequests(ctx, tasks, messages)
}

// Stop gracefully stops the worker
func (w *Worker) Stop() {
	close(w.done)
}

func (w *Worker) ProcessSingleTask(ctx context.Context, out chan Message) (*queue.Task, error) {
	task, err := w.queue.Pop(ctx)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("no tasks available")
	}
	tasks := make(chan *queue.Task)
	messages := make(chan Message)
	go func() {
		for msg := range messages {
			out <- msg
			if msg.Closed || msg.Done {
				close(out)
				close(messages)
				close(tasks)
				return
			}
		}
	}()

	go func() {
		tasks <- task
	}()
	go func() {
		if err := w.processRequests(ctx, tasks, messages); err != nil {
			panic(err)
		}
	}()

	return task, nil
}

func (w *Worker) taskError(ctx context.Context, task *queue.Task, messages chan Message, err error) {
	messages <- Message{
		Type: "error",
		Data: err.Error(),
	}

	if task == nil {
		return
	}

	task.Status = queue.TaskStatusFailed
	task.Error = err.Error()
	task.UpdatedAt = time.Now()
	if err := w.queue.Update(ctx, task); err != nil {
		messages <- Message{
			Type: "error",
			Data: fmt.Sprintf("Failed to update task: %v", err),
		}
	}
}

// func taskMessage(task *queue.Task) Message {
// 	data := make(map[string]string)
// 	data["task"] = task.ID
// 	data["plugin"] = task.PluginID
// 	data["hash"] = task.Hash

// 	msg := ""
// 	for k, v := range data {
// 		msg += fmt.Sprintf("[%s: %s]", k, v)
// 	}

// 	return Message{
// 		Type: "info",
// 		Data: msg,
// 	}
// }

// ProcessSingleTask processes a single task and returns the result
func (w *Worker) processRequests(ctx context.Context, tasks <-chan *queue.Task, messages chan Message) error {
	messages <- Message{
		Type: "info",
		Data: "Processing tasks",
	}

	// Create request channels for each plugin
	pluginReqs := make(map[string]chan *pb.PluginRequest)
	taskmap := make(map[string]map[string]*queue.Task)

	// Start plugin processes
	for id, cfg := range w.plugins {
		reqs := make(chan *pb.PluginRequest)
		pluginReqs[id] = reqs
		taskmap[id] = make(map[string]*queue.Task)

		go func() {
			resps, err := pluginutil.Execute(ctx, cfg, reqs)
			if err != nil {
				messages <- Message{
					Type: "error",
					Data: fmt.Sprintf("Failed to start plugin %s: %v", id, err),
				}
				return
			}

			// Process responses
			for resp := range resps {
				if resp.Closed {
					close(reqs)
					return
				}
				// TODO: consider removing tasks from map when they are done or failed
				task := taskmap[resp.PluginId][resp.RequestId]
				if resp.Error != "" {
					w.taskError(ctx, task, messages, fmt.Errorf("plugin %s: %s", id, resp.Error))
					continue
				}

				if resp.Item != nil {
					// Store the data item
					_, err := w.storeData(resp.Item)
					if err != nil {
						w.taskError(ctx, task, messages, err)
						continue
					}

					messages <- Message{
						Type: "info",
						Data: "[kind: " + resp.Item.Meta.Kind + "] [data: " + resp.Item.Meta.Id + "] [hash: " + resp.Item.Meta.Hash + "]",
					}
				}

				if resp.Action != nil {
					// Create a new task for this item
					newTask := queue.NewTask(*w.plugins[resp.PluginId], resp.Action)
					if err := w.queue.Push(ctx, newTask); err != nil {
						w.taskError(ctx, task, messages, err)
						continue
					}

					messages <- Message{
						Type: "info",
						Data: "[task: " + newTask.ID + "] [plugin: " + newTask.PluginID + "] [action: " + newTask.Action.Name + "]",
					}
				}

				if resp.Done {
					task.Status = queue.TaskStatusComplete
					task.UpdatedAt = time.Now()

					if err := w.queue.Update(ctx, task); err != nil {
						messages <- Message{
							Type: "error",
							Data: fmt.Sprintf("Failed to update task: %v", err),
						}
					}

					messages <- Message{
						Done: true,
					}
				}
			}
		}()
	}

	// Process tasks
	for task := range tasks {
		messages <- Message{
			Type: "info",
			Data: "Processing task: " + task.Key(),
		}

		taskmap[task.PluginID][task.ID] = task
		// Load the data for this task

		var item *pb.DataItem
		var err error
		if task.Hash != "" {
			item, err = w.LoadData(task.Hash)
			if err != nil {
				return fmt.Errorf("failed to load data: %w", err)
			}
		}

		request := task.Request()
		if item != nil {
			request.Item = item
			request.Operation = "transform"
		} else {
			request.Operation = "extract"
		}

		if item == nil {
			messages <- Message{
				Type: "info",
				Data: "[input: nil]",
			}
		} else {
			messages <- Message{
				Type: "info",
				Data: "[input: " + item.Meta.Id + "] [hash: " + item.Meta.Hash + "]",
			}

			messages <- Message{
				Type: "protobuf",
				Data: item.Meta.String(),
			}
		}

		pluginReqs[request.PluginId] <- request
	}

	// Close all plugin request channels
	for _, reqs := range pluginReqs {
		close(reqs)
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

func (w *Worker) LoadData(hash string) (*pb.DataItem, error) {
	filename := filepath.Join(w.dataDir, hash+".dq")

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer f.Close()

	return dq.Read(f)
}
