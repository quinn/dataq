package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type QueueState = map[string]*Task

// FileQueue implements Queue using a file-based storage system
type FileQueue struct {
	dir   string
	mu    sync.Mutex
	state QueueState
}

// newFileQueue creates a new file-based queue
func newFileQueue(dir string) (*FileQueue, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create queue directory: %w", err)
	}

	q := &FileQueue{dir: dir}

	if err := q.loadState(); err != nil {
		return nil, fmt.Errorf("failed to load queue state: %w", err)
	}

	return q, nil
}

func (q *FileQueue) loadState() error {
	f, err := os.Open(filepath.Join(q.dir, "state.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	var state QueueState
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return err
	}

	q.state = state
	return nil
}

func (q *FileQueue) writeState() error {
	f, err := os.Create(filepath.Join(q.dir, "state.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(q.state)
}

// func (q *FileQueue) readTaskFile(filename string) (*TaskMetadata, error) {
// 	f, err := os.Open(filepath.Join(q.dir, filename))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	var task TaskMetadata
// 	data, err := dq.Read(f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	task.Data = data
// 	return &task, nil
// }

// func (q *FileQueue) writeTaskFile(data *pb.DataItem) error {
// 	filename := q.generateFilename(data)
// 	f, err := os.Create(filepath.Join(q.dir, filename))
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	return dq.WriteDataItem(f, data)
// }

func (q *FileQueue) Push(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Update metadata
	q.state[task.ID] = task

	// Update state
	if err := q.writeState(); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}

func (q *FileQueue) Pop(ctx context.Context) (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Find first pending task
	var task *Task
	for _, t := range q.state {
		if t.Status == TaskStatusPending {
			task = t
			break
		}
	}

	if task == nil {
		return nil, nil
	}

	// Update task status
	task.Status = TaskStatusProcessing
	task.UpdatedAt = time.Now()

	// Update state
	if err := q.writeState(); err != nil {
		return nil, fmt.Errorf("failed to write state: %w", err)
	}

	return task, nil
}

func (q *FileQueue) Update(ctx context.Context, meta *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Check if task exists
	if _, ok := q.state[meta.ID]; !ok {
		return fmt.Errorf("task %s not found", meta.ID)
	}

	// Update metadata
	existing := q.state[meta.ID]
	existing.Status = meta.Status
	existing.UpdatedAt = time.Now()

	// Update state
	if err := q.writeState(); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}

func (q *FileQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var tasks []*Task
	for _, task := range q.state {
		if status == "" || task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (q *FileQueue) Close() error {
	return q.writeState()
}
