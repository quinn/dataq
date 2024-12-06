package queue

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type operation struct {
	Type      string    `json:"type"` // "push", "update", "pop"
	Task      *Task     `json:"task"`
	Timestamp time.Time `json:"timestamp"`
}

// FileQueue implements Queue using a file-based storage system
type FileQueue struct {
	dir     string
	mu      sync.Mutex
	state   []string         // ordered list of tasks
	taskMap map[string]*Task // for quick lookups
	file    *os.File
	writer  *bufio.Writer
}

// newFileQueue creates a new file-based queue
func newFileQueue(dir string) (*FileQueue, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create queue directory: %w", err)
	}

	q := &FileQueue{
		dir:     dir,
		taskMap: make(map[string]*Task),
		state:   make([]string, 0),
	}

	if err := q.loadState(); err != nil {
		return nil, fmt.Errorf("failed to load queue state: %w", err)
	}

	return q, nil
}

func (q *FileQueue) loadState() error {
	f, err := os.OpenFile(filepath.Join(q.dir, "queue.log"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open queue log: %w", err)
	}
	q.file = f
	q.writer = bufio.NewWriter(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var op operation
		if err := json.Unmarshal(scanner.Bytes(), &op); err != nil {
			continue // skip corrupted entries
		}

		switch op.Type {
		case "push":
			q.state = append(q.state, op.Task.Key())
			q.taskMap[op.Task.Key()] = op.Task
		case "update":
			if existing := q.taskMap[op.Task.Key()]; existing != nil {
				*existing = *op.Task
			}
		case "pop":
			if len(q.state) > 0 && q.state[0] == op.Task.Key() {
				q.state = q.state[1:]
			}
		}
	}

	return scanner.Err()
}

func (q *FileQueue) writeOperation(op operation) error {
	data, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	if _, err := q.writer.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	if err := q.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	return nil
}

func (q *FileQueue) Push(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	op := operation{
		Type:      "push",
		Task:      task,
		Timestamp: time.Now(),
	}

	if err := q.writeOperation(op); err != nil {
		return fmt.Errorf("failed to write push operation: %w", err)
	}

	q.state = append(q.state, task.Key())
	q.taskMap[task.Key()] = task
	return nil
}

func (q *FileQueue) Pop(ctx context.Context) (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Find first pending task
	var taskIndex int = -1
	var task *Task
	for i, tid := range q.state {
		t := q.taskMap[tid]
		if t.Status == TaskStatusPending {
			task = t
			taskIndex = i
			break
		}
	}

	if task == nil {
		return nil, nil
	}

	// Update task status
	task.Status = TaskStatusProcessing
	task.UpdatedAt = time.Now()

	op := operation{
		Type:      "pop",
		Task:      task,
		Timestamp: time.Now(),
	}

	if err := q.writeOperation(op); err != nil {
		return nil, fmt.Errorf("failed to write pop operation: %w", err)
	}

	// Remove from state
	q.state = append(q.state[:taskIndex], q.state[taskIndex+1:]...)

	return task, nil
}

func (q *FileQueue) Update(ctx context.Context, task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	existing := q.taskMap[task.Key()]
	if existing == nil {
		return fmt.Errorf("task not found: %s", task.Key())
	}

	op := operation{
		Type:      "update",
		Task:      task,
		Timestamp: time.Now(),
	}

	if err := q.writeOperation(op); err != nil {
		return fmt.Errorf("failed to write update operation: %w", err)
	}

	*existing = *task
	return nil
}

func (q *FileQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var tasks []*Task
	for _, tid := range q.state {
		t := q.taskMap[tid]
		if status == "" || t.Status == status {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (q *FileQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.writer != nil {
		if err := q.writer.Flush(); err != nil {
			return err
		}
	}
	if q.file != nil {
		return q.file.Close()
	}
	return nil
}
