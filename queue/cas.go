package queue

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.quinn.io/dataq/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ClaimsQueue struct {
	db *gorm.DB
}

// NewSQLiteQueue creates a new SQLite-backed queue
func NewClaimsQueue(conn *sql.DB) (*SQLiteQueue, error) {
	// Convert sql.DB to gorm.DB
	db, err := gorm.Open(sqlite.Dialector{Conn: conn})
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&TaskModel{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &SQLiteQueue{db: db}, nil
}

func (q *SQLiteQueue) Push(ctx context.Context, task *schema.Task) error {
	model, err := TaskModelFromTask(task)
	if err != nil {
		return fmt.Errorf("failed to convert task to model: %w", err)
	}

	if err := q.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

func (q *SQLiteQueue) Pop(ctx context.Context) (*schema.Task, error) {
	var model TaskModel

	err := q.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the first pending task
		if err := tx.Where("status = ?", TaskStatusPending).First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return fmt.Errorf("failed to find pending task: %w", err)
		}

		// Update task status to processing
		model.Status = string(TaskStatusProcessing)
		model.UpdatedAt = time.Now()
		if err := tx.Save(&model).Error; err != nil {
			return fmt.Errorf("failed to update task status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if model.ID == "" { // No task found
		return nil, nil
	}

	return model.ToTask()
}

func (q *SQLiteQueue) Update(ctx context.Context, task *schema.Task) error {
	model, err := TaskModelFromTask(task)
	if err != nil {
		return fmt.Errorf("failed to convert task to model: %w", err)
	}

	result := q.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found: %s", task.Uid)
	}

	return nil
}

func (q *SQLiteQueue) List(ctx context.Context, status TaskStatus) ([]*schema.Task, error) {
	var models []TaskModel
	if err := q.db.WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	tasks := make([]*schema.Task, 0, len(models))
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
