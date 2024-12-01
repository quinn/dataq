package queue

import (
	"context"
	"fmt"
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
type TaskMetadata struct {
	ID        string
	Status    TaskStatus
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Task struct {
	Meta TaskMetadata
	Data *pb.DataItem
}

func NewTask(data *pb.DataItem) *Task {
	return &Task{
		Meta: TaskMetadata{
			ID:        fmt.Sprintf("%s_%s", data.Meta.PluginId, data.Meta.Id),
			Status:    TaskStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Data: data,
	}
}

func InitialTask(pluginID string) *Task {
	return &Task{
		Meta: TaskMetadata{
			ID:        fmt.Sprintf("%s_initial", pluginID),
			Status:    TaskStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
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
	// Push adds a task to the queue
	Push(ctx context.Context, task *Task) error

	// Pop removes and returns the next task from the queue
	Pop(ctx context.Context) (*Task, error)

	// Update updates an existing task in the queue
	Update(ctx context.Context, meta *TaskMetadata) error

	// Close closes the queue and releases any resources
	Close() error

	// List returns all tasks in the queue, optionally filtered by status
	List(ctx context.Context, status TaskStatus) ([]*TaskMetadata, error)
}
