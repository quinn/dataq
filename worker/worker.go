package worker

import (
	"context"
	"log"
	"os"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/queue"
)

type Worker struct {
	queue   queue.Queue
	cas     cas.Storage
	plugins map[string]*config.Plugin
	dataDir string
	done    chan struct{}
}

func New(q queue.Queue, plugins []*config.Plugin, dataDir string, c cas.Storage) *Worker {
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
		cas:     c,
	}
}

func (w *Worker) Start(ctx context.Context, messages chan Message) error {
	sendInfo(messages, "Starting task processing loop")

	tasks := make(chan *queue.Task)
	go w.dequeueTasks(ctx, tasks, messages)

	return w.runTaskProcessingLoop(ctx, tasks, messages)
}

func (w *Worker) Stop() {
	close(w.done)
}

func (w *Worker) dequeueTasks(ctx context.Context, tasks chan<- *queue.Task, messages chan Message) {
	defer close(tasks)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		task, err := w.queue.Pop(ctx)
		if err != nil {
			sendErrorf(messages, "Failed to pop task: %v", err)
			return
		}

		if task == nil {
			continue
		}

		tasks <- task
	}
}

// ProcessSingleTask is similar to the original approach, but simplified.
// It's basically a one-off run of the whole pipeline for a single task.
func (w *Worker) ProcessSingleTask(ctx context.Context, out chan Message) (*queue.Task, error) {
	task, err := w.queue.Pop(ctx)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, err
	}

	// For single task processing, we create a tasks chan just for it.
	tasks := make(chan *queue.Task)
	go func() {
		tasks <- task
		close(tasks)
	}()

	messages := make(chan Message)
	go func() {
		for msg := range messages {
			out <- msg
			if msg.Closed || msg.Done {
				close(out)
				return
			}
		}
	}()

	go func() {
		if err := w.runTaskProcessingLoop(ctx, tasks, messages); err != nil {
			// If we panic here, it's the same behavior as original code
			panic(err)
		}
	}()

	return task, nil
}
