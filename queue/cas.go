package queue

import (
	"context"
	"fmt"
	"time"

	"go.quinn.io/dataq/claims"
	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/schema"
	"gorm.io/gorm"
)

type ClaimsQueue struct {
	index *index.ClaimsIndexer
}

// NewClaimsQueue creates a queue stored in cas and indexed in a database
func NewClaimsQueue(index *index.ClaimsIndexer) *ClaimsQueue {
	return &ClaimsQueue{index: index}
}

func (q *ClaimsQueue) Push(ctx context.Context, task *schema.Task) error {
	claim := &claims.Claim{
		Kind:    "task",
		Payload: task,
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
