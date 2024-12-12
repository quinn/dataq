package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	pb "go.quinn.io/dataq/proto"
)

// TaskModel represents the database model for tasks
type TaskModel struct {
	ID           string `gorm:"primaryKey"`
	PluginID     string
	Status       string
	Error        string
	Hash         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PluginConfig string // JSON encoded map
	Action       string // JSON encoded pb.Action
}

// ToTask converts the database model to a Task
func (m *TaskModel) ToTask() (*Task, error) {
	var pluginConfig map[string]string
	if err := json.Unmarshal([]byte(m.PluginConfig), &pluginConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin config: %w", err)
	}

	var action pb.Action
	if err := json.Unmarshal([]byte(m.Action), &action); err != nil {
		return nil, fmt.Errorf("failed to unmarshal action: %w", err)
	}

	return &Task{
		ID:           m.ID,
		PluginID:     m.PluginID,
		Status:       TaskStatus(m.Status),
		Error:        m.Error,
		Hash:         m.Hash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		PluginConfig: pluginConfig,
		Action:       &action,
	}, nil
}

// FromTask converts a Task to a database model
func TaskModelFromTask(t *Task) (*TaskModel, error) {
	pluginConfig, err := json.Marshal(t.PluginConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal plugin config: %w", err)
	}

	action, err := json.Marshal(t.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal action: %w", err)
	}

	return &TaskModel{
		ID:           t.ID,
		PluginID:     t.PluginID,
		Status:       string(t.Status),
		Error:        t.Error,
		Hash:         t.Hash,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		PluginConfig: string(pluginConfig),
		Action:       string(action),
	}, nil
}

type SQLiteQueue struct {
	db *gorm.DB
}

// NewSQLiteQueue creates a new SQLite-backed queue
func NewSQLiteQueue(path string) (*SQLiteQueue, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&TaskModel{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &SQLiteQueue{db: db}, nil
}

func (q *SQLiteQueue) Push(ctx context.Context, task *Task) error {
	model, err := TaskModelFromTask(task)
	if err != nil {
		return fmt.Errorf("failed to convert task to model: %w", err)
	}

	if err := q.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

func (q *SQLiteQueue) Pop(ctx context.Context) (*Task, error) {
	var model TaskModel

	err := q.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the first pending task
		if err := tx.Where("status = ?", TaskStatusPending).First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return fmt.Errorf("failed to find pending task: %w", err)
		}

		// Delete the task
		if err := tx.Delete(&model).Error; err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		return nil
	})

	if err == nil {
		return model.ToTask()
	}
	return nil, err
}

func (q *SQLiteQueue) Update(ctx context.Context, task *Task) error {
	model, err := TaskModelFromTask(task)
	if err != nil {
		return fmt.Errorf("failed to convert task to model: %w", err)
	}

	result := q.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found: %s", task.ID)
	}

	return nil
}

func (q *SQLiteQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	var models []TaskModel
	if err := q.db.WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	tasks := make([]*Task, 0, len(models))
	for _, model := range models {
		task, err := model.ToTask()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (q *SQLiteQueue) Close() error {
	// don't close db since it is passed in
	return nil
}
