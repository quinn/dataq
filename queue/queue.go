package queue

import (
	"context"
	"fmt"
	"time"

	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/hash"
	"go.quinn.io/dataq/schema"
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

// func (t *schema.Task) Key() string {
// 	return t.Uid
// }

// func (t *schema.Task) Request() *schema.PluginRequest {
// 	return &schema.PluginRequest{
// 		PluginId:     t.PluginID,
// 		Id:           t.ID,
// 		Operation:    "extract",
// 		PluginConfig: t.PluginConfig,
// 		Action:       t.Action,
// 	}
// }

// NewTask creates a new TaskMetadata instance
func NewTask(plugin config.Plugin, action *schema.Action) *schema.Task {
	return &schema.Task{
		Uid:          hash.UID(),
		Status:       string(TaskStatusPending),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		ActionUid:    action.Id,
		PluginConfig: plugin.Config,
		PluginId:     plugin.ID,
	}
}

func NewExtractTask(plugin config.Plugin, item *schema.DataItem) *schema.Task {
	return &schema.Task{
		Uid:          hash.UID(),
		Status:       string(TaskStatusPending),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		PluginConfig: plugin.Config,
		PluginId:     plugin.ID,
		DataItemUid:  item.Meta.Id,
	}
}

// InitialTask creates an initial task metadata for a plugin
func InitialTask(plugin config.Plugin) *schema.Task {
	return NewTask(plugin, &schema.Action{
		Id:   hash.UID(),
		Name: "initial",
	})
}

// Queue defines the interface for task queues
type Queue interface {
	// Push adds a task metadata to the queue
	Push(ctx context.Context, meta *schema.Task) error

	// Pop removes and returns the next task metadata from the queue
	Pop(ctx context.Context) (*schema.Task, error)

	// Update updates an existing task metadata in the queue
	Update(ctx context.Context, meta *schema.Task) error

	// Close closes the queue and releases any resources
	Close() error

	// List returns all task metadata in the queue, optionally filtered by status
	List(ctx context.Context, status TaskStatus) ([]*schema.Task, error)
}
