package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

var (
	queueBucket = []byte("queue")
)

// BoltQueue implements Queue using BBolt
type BoltQueue struct {
	db *bbolt.DB
}

// NewBoltQueue creates a new BoltQueue instance
func NewBoltQueue(opts ...Option) (Queue, error) {
	options := &Options{
		Path: "queue.db",
	}

	for _, opt := range opts {
		opt(options)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(options.Path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create queue directory: %v", err)
	}

	// Open BBolt database
	db, err := bbolt.Open(options.Path, 0600, &bbolt.Options{
		Timeout:      5 * time.Second,
		NoSync:       false,
		NoGrowSync:   false,
		FreelistType: bbolt.FreelistArrayType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open queue database: %v", err)
	}

	// Create bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(queueBucket)
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create queue bucket: %v", err)
	}

	return &BoltQueue{db: db}, nil
}

// Push adds a task to the queue
func (q *BoltQueue) Push(ctx context.Context, task *Task) error {
	return q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(queueBucket)

		// Marshal task to JSON
		data, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %v", err)
		}

		// Store task
		return b.Put([]byte(task.ID), data)
	})
}

// Pop removes and returns the next task from the queue
func (q *BoltQueue) Pop(ctx context.Context) (*Task, error) {
	var task *Task

	err := q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(queueBucket)
		c := b.Cursor()

		// Get first pending task
		k, v := c.First()
		for ; k != nil; k, v = c.Next() {
			var t Task
			if err := json.Unmarshal(v, &t); err != nil {
				continue
			}

			if t.Status == TaskStatusPending {
				task = &t
				return b.Delete(k)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pop task: %v", err)
	}

	return task, nil
}

// Get returns a task by ID
func (q *BoltQueue) Get(ctx context.Context, id string) (*Task, error) {
	var task *Task

	err := q.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(queueBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("task not found: %s", id)
		}

		task = &Task{}
		if err := json.Unmarshal(v, task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return task, nil
}

// Update updates a task in the queue
func (q *BoltQueue) Update(ctx context.Context, task *Task) error {
	return q.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(queueBucket)

		// Marshal task to JSON
		data, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %v", err)
		}

		// Store task
		return b.Put([]byte(task.ID), data)
	})
}

// List returns all tasks in the queue, optionally filtered by status
func (q *BoltQueue) List(ctx context.Context, status *TaskStatus) ([]*Task, error) {
	var tasks []*Task

	err := q.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(queueBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				continue
			}
			
			if status != nil && task.Status != *status {
				continue
			}
			
			tasks = append(tasks, &task)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %v", err)
	}

	return tasks, nil
}

// Close closes the queue
func (q *BoltQueue) Close() error {
	return q.db.Close()
}
