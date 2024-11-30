package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.quinn.io/dataq/dq"
)

// FileQueue implements Queue using a file-based storage system
type FileQueue struct {
	dir      string
	mu       sync.Mutex
	metadata map[string]*Task // in-memory index of tasks
}

// NewFileQueue creates a new file-based queue
func NewFileQueue(dir string) (*FileQueue, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create queue directory: %w", err)
	}

	q := &FileQueue{
		dir:      dir,
		metadata: make(map[string]*Task),
	}

	// Load existing tasks
	if err := q.loadExistingTasks(); err != nil {
		return nil, fmt.Errorf("failed to load existing tasks: %w", err)
	}

	return q, nil
}

func (q *FileQueue) loadExistingTasks() error {
	entries, err := os.ReadDir(q.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".task" {
			continue
		}

		task, err := q.readTaskFile(entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read task %s: %w", entry.Name(), err)
		}

		q.metadata[task.ID] = task
	}

	return nil
}

func (q *FileQueue) readTaskFile(filename string) (*Task, error) {
	f, err := os.Open(filepath.Join(q.dir, filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var task Task
	if err := json.NewDecoder(f).Decode(&task); err != nil {
		return nil, err
	}

	// If task has data, read it using dq format
	if task.Data != nil {
		dataFile := filepath.Join(q.dir, task.ID+".data")
		df, err := os.Open(dataFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil {
			task.Data, err = dq.Read(df)
			df.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read data: %w", err)
			}
		}
	}

	return &task, nil
}

func (q *FileQueue) writeTaskFile(task *Task) error {
	// Write task metadata
	f, err := os.Create(filepath.Join(q.dir, task.ID+".task"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(task); err != nil {
		return err
	}

	// If task has data, write it using dq format
	if task.Data != nil {
		df, err := os.Create(filepath.Join(q.dir, task.ID+".data"))
		if err != nil {
			return err
		}
		defer df.Close()

		if err := dq.WriteDataItem(df, task.Data); err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}
	}

	return nil
}

func (q *FileQueue) Push(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if task.ID == "" {
		return fmt.Errorf("task ID is required")
	}

	if _, exists := q.metadata[task.ID]; exists {
		return fmt.Errorf("task %s already exists", task.ID)
	}

	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	task.UpdatedAt = task.CreatedAt

	if err := q.writeTaskFile(task); err != nil {
		return fmt.Errorf("failed to write task: %w", err)
	}

	q.metadata[task.ID] = task
	return nil
}

func (q *FileQueue) Pop(ctx context.Context) (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var oldestID string
	var oldestTime time.Time

	// Find the oldest pending task
	for id, task := range q.metadata {
		if task.Status != TaskStatusPending {
			continue
		}
		if oldestID == "" || task.CreatedAt.Before(oldestTime) {
			oldestID = id
			oldestTime = task.CreatedAt
		}
	}

	if oldestID == "" {
		return nil, fmt.Errorf("no pending tasks")
	}

	task := q.metadata[oldestID]
	delete(q.metadata, oldestID)

	// Remove task files
	os.Remove(filepath.Join(q.dir, task.ID+".task"))
	os.Remove(filepath.Join(q.dir, task.ID+".data"))

	return task, nil
}

func (q *FileQueue) Update(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.metadata[task.ID]; !exists {
		return fmt.Errorf("task %s not found", task.ID)
	}

	task.UpdatedAt = time.Now()

	if err := q.writeTaskFile(task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	q.metadata[task.ID] = task
	return nil
}

func (q *FileQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var tasks []*Task
	for _, task := range q.metadata {
		if status == "" || task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (q *FileQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Clear the metadata map
	q.metadata = nil
	return nil
}
