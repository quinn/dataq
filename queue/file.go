package queue

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// QueueState represents the current state of tasks in the queue
type QueueState struct {
	Tasks map[string]TaskStatus `json:"tasks"` // map of task ID to status
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

	// Load existing tasks and state
	if err := q.loadState(); err != nil {
		return nil, fmt.Errorf("failed to load queue state: %w", err)
	}

	return q, nil
}

func (q *FileQueue) loadState() error {
	// Read state file
	state, err := q.readState()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if state == nil {
		state = &QueueState{Tasks: make(map[string]TaskStatus)}
	}

	// Load tasks from state
	entries, err := os.ReadDir(q.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".dq" {
			continue
		}

		task, err := q.readTaskFile(entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read task %s: %w", entry.Name(), err)
		}

		// Update task status from state
		if status, ok := state.Tasks[task.ID]; ok {
			task.Status = status
		}

		q.metadata[task.ID] = task
	}

	return nil
}

func (q *FileQueue) readState() (*QueueState, error) {
	f, err := os.Open(filepath.Join(q.dir, "state.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var state QueueState
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
}

func (q *FileQueue) writeState() error {
	state := QueueState{
		Tasks: make(map[string]TaskStatus),
	}

	for id, task := range q.metadata {
		state.Tasks[id] = task.Status
	}

	f, err := os.Create(filepath.Join(q.dir, "state.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(state)
}

func (q *FileQueue) generateFilename(task *Task) string {
	// DataItem must exist since this is only called for plugin response items
	h := sha256.New()
	h.Write(task.Data.RawData)
	return hex.EncodeToString(h.Sum(nil)) + ".dq"
}

func (q *FileQueue) readTaskFile(filename string) (*Task, error) {
	f, err := os.Open(filepath.Join(q.dir, filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var task Task
	data, err := dq.Read(f)
	if err != nil {
		return nil, err
	}

	task.Data = data
	return &task, nil
}

func (q *FileQueue) writeTaskFile(task *Task) error {
	filename := q.generateFilename(task)
	f, err := os.Create(filepath.Join(q.dir, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	return dq.Write(f, task, task.Data.RawData)
}

func (q *FileQueue) Push(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Only write file for tasks with data (from plugin responses)
	if task.Data != nil {
		if err := q.writeTaskFile(task); err != nil {
			return fmt.Errorf("failed to write task file: %w", err)
		}
	}

	// Update metadata
	q.metadata[task.ID] = task

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
	for _, t := range q.metadata {
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

func (q *FileQueue) Update(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Check if task exists
	if _, ok := q.metadata[task.ID]; !ok {
		return fmt.Errorf("task %s not found", task.ID)
	}

	// Update metadata
	q.metadata[task.ID] = task

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
	for _, task := range q.metadata {
		if status == "" || task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (q *FileQueue) Close() error {
	return q.writeState()
}
