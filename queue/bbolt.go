package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	"go.quinn.io/dataq/dq"
)

var (
	tasksBucket = []byte("tasks")
)

type BoltQueue struct {
	db *bbolt.DB
}

func NewBoltQueue(path string) (*BoltQueue, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tasksBucket)
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &BoltQueue{db: db}, nil
}

func (q *BoltQueue) Push(ctx context.Context, task *Task) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	task.UpdatedAt = time.Now()

	// Serialize task metadata
	metadata := struct {
		ID        string            `json:"id"`
		PluginID  string            `json:"plugin_id"`
		Config    map[string]string `json:"config"`
		Status    TaskStatus        `json:"status"`
		Error     string            `json:"error"`
		CreatedAt time.Time         `json:"created_at"`
		UpdatedAt time.Time         `json:"updated_at"`
	}{
		ID:        task.ID,
		PluginID:  task.PluginID,
		Config:    task.Config,
		Status:    task.Status,
		Error:     task.Error,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	// Combine metadata and data
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	value := append(metadataJSON, data...)

	err = q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		return b.Put([]byte(task.ID), value)
	})
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

func (q *BoltQueue) Pop(ctx context.Context) (*Task, error) {
	var task *Task

	err := q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var metadata struct {
				ID        string            `json:"id"`
				PluginID  string            `json:"plugin_id"`
				Config    map[string]string `json:"config"`
				Status    TaskStatus        `json:"status"`
				Error     string            `json:"error"`
				CreatedAt time.Time         `json:"created_at"`
				UpdatedAt time.Time         `json:"updated_at"`
			}

			// Find metadata length by looking for first non-JSON byte
			var metadataLen int
			for i := 0; i < len(v); i++ {
				if !json.Valid(v[:i+1]) {
					metadataLen = i
					break
				}
			}

			if err := json.Unmarshal(v[:metadataLen], &metadata); err != nil {
				continue
			}

			if metadata.Status == TaskStatusPending {
				task = &Task{
					ID:        metadata.ID,
					PluginID:  metadata.PluginID,
					Config:    metadata.Config,
					Status:    metadata.Status,
					Error:     metadata.Error,
					CreatedAt: metadata.CreatedAt,
					UpdatedAt: metadata.UpdatedAt,
				}

				// Deserialize data using dq package
				if len(v) > metadataLen {
					data, err := dq.Read(bytes.NewReader(v[metadataLen:]))
					if err != nil {
						return fmt.Errorf("failed to deserialize data: %w", err)
					}
					task.Data = data
				}

				if err := b.Delete(k); err != nil {
					return fmt.Errorf("failed to delete task: %w", err)
				}
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pop task: %w", err)
	}

	return task, nil
}

func (q *BoltQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	var tasks []*Task

	err := q.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var metadata struct {
				ID        string            `json:"id"`
				PluginID  string            `json:"plugin_id"`
				Config    map[string]string `json:"config"`
				Status    TaskStatus        `json:"status"`
				Error     string            `json:"error"`
				CreatedAt time.Time         `json:"created_at"`
				UpdatedAt time.Time         `json:"updated_at"`
			}

			// Find metadata length by looking for first non-JSON byte
			var metadataLen int
			for i := 0; i < len(v); i++ {
				if !json.Valid(v[:i+1]) {
					metadataLen = i
					break
				}
			}

			if err := json.Unmarshal(v[:metadataLen], &metadata); err != nil {
				continue
			}

			if status == "" || metadata.Status == status {
				task := &Task{
					ID:        metadata.ID,
					PluginID:  metadata.PluginID,
					Config:    metadata.Config,
					Status:    metadata.Status,
					Error:     metadata.Error,
					CreatedAt: metadata.CreatedAt,
					UpdatedAt: metadata.UpdatedAt,
				}

				// Deserialize data using dq package
				if len(v) > metadataLen {
					data, err := dq.Read(bytes.NewReader(v[metadataLen:]))
					if err != nil {
						return fmt.Errorf("failed to deserialize data: %w", err)
					}
					task.Data = data
				}

				tasks = append(tasks, task)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (q *BoltQueue) Update(ctx context.Context, task *Task) error {
	task.UpdatedAt = time.Now()

	// Serialize task metadata
	metadata := struct {
		ID        string            `json:"id"`
		PluginID  string            `json:"plugin_id"`
		Config    map[string]string `json:"config"`
		Status    TaskStatus        `json:"status"`
		Error     string            `json:"error"`
		CreatedAt time.Time         `json:"created_at"`
		UpdatedAt time.Time         `json:"updated_at"`
	}{
		ID:        task.ID,
		PluginID:  task.PluginID,
		Config:    task.Config,
		Status:    task.Status,
		Error:     task.Error,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	// Combine metadata and data
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	value := append(metadataJSON, data...)

	err = q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		return b.Put([]byte(task.ID), value)
	})
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (q *BoltQueue) Close() error {
	return q.db.Close()
}
