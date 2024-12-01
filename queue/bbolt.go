package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	tasksBucket = []byte("tasks")
	Delimiter   = []byte("\n")
)

type BoltQueue struct {
	db *bbolt.DB
}

func newBoltQueue(path string) (*BoltQueue, error) {
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

func (q *BoltQueue) Push(ctx context.Context, task *TaskMetadata) error {
	return q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)

		// Serialize metadata
		meta, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		return b.Put([]byte(task.ID), meta)
	})
}

func (q *BoltQueue) Pop(ctx context.Context) (*TaskMetadata, error) {
	var task *TaskMetadata

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
		task = &TaskMetadata{}

		// Find metadata JSON end
		var i int
		for i = 0; i < len(v); i++ {
			if v[i] == '{' {
				break
			}
		}
		metadataJSON := v[:i]

		// Parse metadata
		if err := json.Unmarshal(metadataJSON, &task); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

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

func (q *BoltQueue) List(ctx context.Context, status TaskStatus) ([]*TaskMetadata, error) {
	var tasks []*TaskMetadata

	err := q.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)

		return b.ForEach(func(k, v []byte) error {
			// Find delimiter between metadata and data
			i := bytes.Index(v, Delimiter)
			if i == -1 {
				return fmt.Errorf("invalid task format: no delimiter found")
			}

			// Parse metadata
			var meta TaskMetadata
			if err := json.Unmarshal(v[:i], &meta); err != nil {
				return fmt.Errorf("failed to unmarshal metadata: %w", err)
			}

			// Only include tasks matching the requested status
			if status == "" || meta.Status == status {
				tasks = append(tasks, &meta)
			}

			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (q *BoltQueue) Update(ctx context.Context, meta *TaskMetadata) error {
	return q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tasksBucket)

		// Get existing value
		v := b.Get([]byte(meta.ID))
		if v == nil {
			return fmt.Errorf("task not found: %s", meta.ID)
		}

		// Find delimiter between metadata and data
		i := bytes.Index(v, Delimiter)
		if i == -1 {
			return fmt.Errorf("invalid task format: no delimiter found")
		}

		// Serialize new metadata
		metadataJSON, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		// Combine new metadata with existing data
		value := append(metadataJSON, Delimiter...)
		value = append(value, v[i+len(Delimiter):]...)

		return b.Put([]byte(meta.ID), value)
	})
}

func (q *BoltQueue) Close() error {
	return q.db.Close()
}
