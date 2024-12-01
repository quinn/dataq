package queue

import (
	"context"
	"time"

	pb "go.quinn.io/dataq/proto"
)

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusComplete   TaskStatus = "complete"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// Task represents a unit of work to be processed
type Task struct {
	ID        string
	PluginID  string
	Config    map[string]string
	Data      *pb.DataItem
	Status    TaskStatus
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Queue defines the interface for task queues
type Queue interface {
	// Push adds a task to the queue
	Push(ctx context.Context, task *Task) error

	// Pop removes and returns the next task from the queue
	Pop(ctx context.Context) (*Task, error)

	// Update updates an existing task in the queue
	Update(ctx context.Context, task *Task) error

	// Close closes the queue and releases any resources
	Close() error

	// List returns all tasks in the queue, optionally filtered by status
	List(ctx context.Context, status TaskStatus) ([]*Task, error)
}
