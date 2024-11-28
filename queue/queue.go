package queue

import (
	"context"
	"time"

	pb "go.quinn.io/dataq/proto"
)

// TaskStatus represents the current state of a task
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusRunning
	TaskStatusCompleted
	TaskStatusFailed
)

// Task represents a unit of work to be processed
type Task struct {
	ID        string
	PluginID  string
	Config    map[string]string
	Status    TaskStatus
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Result    *pb.PluginResponse
}

// Queue defines the interface for queue implementations
type Queue interface {
	// Push adds a new task to the queue
	Push(ctx context.Context, task *Task) error

	// Pop removes and returns the next task from the queue
	// If no task is available, returns nil and no error
	Pop(ctx context.Context) (*Task, error)

	// Get retrieves a task by ID without removing it
	Get(ctx context.Context, id string) (*Task, error)

	// Update updates the status and result of a task
	Update(ctx context.Context, id string, status TaskStatus, result *pb.PluginResponse, err string) error

	// List returns all tasks matching the given status
	// If status is nil, returns all tasks
	List(ctx context.Context, status *TaskStatus) ([]*Task, error)

	// Close closes the queue and releases any resources
	Close() error
}

// QueueOption represents an option when creating a new queue
type QueueOption func(*QueueOptions)

// QueueOptions contains all queue configuration options
type QueueOptions struct {
	Path string // Path to the queue storage directory
}

// WithPath sets the storage path for the queue
func WithPath(path string) QueueOption {
	return func(o *QueueOptions) {
		o.Path = path
	}
}
