package queue

import (
	"context"
	"fmt"
	"time"
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

// Task represents metadata about a unit of work to be processed
type TaskMetadata struct {
	ID        string
	Status    TaskStatus
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DataHash  string // SHA-256 hash of the data, used to reference the actual data
}

// NewTaskMetadata creates a new TaskMetadata instance
func NewTaskMetadata(pluginID, itemID string, hash string) *TaskMetadata {
	return &TaskMetadata{
		ID:        fmt.Sprintf("%s_%s", pluginID, itemID),
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DataHash:  hash,
	}
}

// InitialTask creates an initial task metadata for a plugin
func InitialTask(pluginID string) *TaskMetadata {
	return &TaskMetadata{
		ID:        fmt.Sprintf("%s_initial", pluginID),
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewQueue(queueType, path string) (Queue, error) {
	switch queueType {
	case "sqlite":
		return newSQLiteQueue(path)
	case "bbolt":
		return newBoltQueue(path)
	case "file":
		return newFileQueue(path)
	default:
		return nil, fmt.Errorf("unknown queue type: %s", queueType)
	}
}

// Queue defines the interface for task queues
type Queue interface {
	// Push adds a task metadata to the queue
	Push(ctx context.Context, meta *TaskMetadata) error

	// Pop removes and returns the next task metadata from the queue
	Pop(ctx context.Context) (*TaskMetadata, error)

	// Update updates an existing task metadata in the queue
	Update(ctx context.Context, meta *TaskMetadata) error

	// Close closes the queue and releases any resources
	Close() error

	// List returns all task metadata in the queue, optionally filtered by status
	List(ctx context.Context, status TaskStatus) ([]*TaskMetadata, error)
}
