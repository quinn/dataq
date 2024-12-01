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
	if task.Meta.ID == "" {
		task.Meta.ID = uuid.New().String()
	}
	if task.Meta.CreatedAt.IsZero() {
		task.Meta.CreatedAt = time.Now()
	}
	task.Meta.UpdatedAt = time.Now()

	// Serialize task metadata
	metadataJSON, err := json.Marshal(task.Meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	// Combine metadata and data
	value := append(metadataJSON, data...)

	err = q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		return b.Put([]byte(task.Meta.ID), value)
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

		// Get first key/value pair
		k, v := c.First()
		if k == nil {
			// No tasks in queue
			return nil
		}

		// Parse task
		task = &Task{
			Meta: TaskMetadata{},
		}

		// Find metadata JSON end
		var i int
		for i = 0; i < len(v); i++ {
			if v[i] == '{' {
				break
			}
		}
		metadataJSON := v[:i]

		// Parse metadata
		if err := json.Unmarshal(metadataJSON, &task.Meta); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		// Parse data
		data := v[i:]
		var buf bytes.Buffer
		if _, err := buf.Write(data); err != nil {
			return fmt.Errorf("failed to write data to buffer: %w", err)
		}

		dataItem, err := dq.ReadDataItem(&buf)
		if err != nil {
			return fmt.Errorf("failed to deserialize data: %w", err)
		}
		task.Data = dataItem

		// Delete task from queue
		return b.Delete(k)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pop task: %w", err)
	}

	if task == nil {
		return nil, nil
	}

	return task, nil
}

func (q *BoltQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	var tasks []*Task

	err := q.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := &Task{
				Meta: TaskMetadata{},
			}

			// Find metadata JSON end
			var i int
			for i = 0; i < len(v); i++ {
				if v[i] == '{' {
					break
				}
			}
			metadataJSON := v[:i]

			// Parse metadata
			if err := json.Unmarshal(metadataJSON, &task.Meta); err != nil {
				return fmt.Errorf("failed to unmarshal metadata: %w", err)
			}

			if status != "" && task.Meta.Status != status {
				continue
			}

			// Parse data
			data := v[i:]
			var buf bytes.Buffer
			if _, err := buf.Write(data); err != nil {
				return fmt.Errorf("failed to write data to buffer: %w", err)
			}

			dataItem, err := dq.ReadDataItem(&buf)
			if err != nil {
				return fmt.Errorf("failed to deserialize data: %w", err)
			}
			task.Data = dataItem

			tasks = append(tasks, task)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (q *BoltQueue) Update(ctx context.Context, task *Task) error {
	task.Meta.UpdatedAt = time.Now()

	// Serialize task metadata
	metadataJSON, err := json.Marshal(task.Meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	// Combine metadata and data
	value := append(metadataJSON, data...)

	err = q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		return b.Put([]byte(task.Meta.ID), value)
	})
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (q *BoltQueue) Close() error {
	return q.db.Close()
}
