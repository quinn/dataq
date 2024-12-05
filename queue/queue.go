package queue

import (
	"context"
	"fmt"
	"time"

	"go.quinn.io/dataq/config"
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

// Task represents metadata about a unit of work to be processed
type Task struct {
	Status    TaskStatus
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Hash      string // SHA-256 hash of the data, used to reference the actual data
	Config    map[string]string
	PluginID  string
	ID        string
}

func (t *Task) Key() string {
	return t.PluginID + "-" + t.ID
}

func (t *Task) Request() *pb.PluginRequest {
	return &pb.PluginRequest{
		PluginId: t.PluginID,
		Id:       t.ID,
		Config:   t.Config,
	}
}

// NewTask creates a new TaskMetadata instance
func NewTask(plugin config.Plugin, itemID string, hash string) *Task {
	return &Task{
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Hash:      hash,
		Config:    plugin.Config,
		PluginID:  plugin.ID,
		ID:        itemID,
	}
}

// InitialTask creates an initial task metadata for a plugin
func InitialTask(plugin config.Plugin) *Task {
	return &Task{
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config:    plugin.Config,
		PluginID:  plugin.ID,
		ID:        "initial",
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
	Push(ctx context.Context, meta *Task) error

	// Pop removes and returns the next task metadata from the queue
	Pop(ctx context.Context) (*Task, error)

	// Update updates an existing task metadata in the queue
	Update(ctx context.Context, meta *Task) error

	// Close closes the queue and releases any resources
	Close() error

	// List returns all task metadata in the queue, optionally filtered by status
	List(ctx context.Context, status TaskStatus) ([]*Task, error)
}
