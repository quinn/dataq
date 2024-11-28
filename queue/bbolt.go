package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
	pb "go.quinn.io/dataq/proto"
)

var (
	// Bucket names
	tasksBucket = []byte("tasks")
	queueBucket = []byte("queue")
)

// BoltQueue implements Queue using BBolt
type BoltQueue struct {
	db *bbolt.DB
}

// NewBoltQueue creates a new BBolt-backed queue
func NewBoltQueue(opts ...QueueOption) (Queue, error) {
	options := &QueueOptions{
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
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open queue database: %v", err)
	}

	// Create buckets if they don't exist
	if err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(tasksBucket); err != nil {
			return fmt.Errorf("failed to create tasks bucket: %v", err)
		}
		if _, err := tx.CreateBucketIfNotExists(queueBucket); err != nil {
			return fmt.Errorf("failed to create queue bucket: %v", err)
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, err
	}

	return &BoltQueue{db: db}, nil
}

func (q *BoltQueue) Push(ctx context.Context, task *Task) error {
	return q.db.Update(func(tx *bbolt.Tx) error {
		// Store task data
		tasksBucket := tx.Bucket(tasksBucket)
		taskData, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %v", err)
		}
		if err := tasksBucket.Put([]byte(task.ID), taskData); err != nil {
			return fmt.Errorf("failed to store task: %v", err)
		}

		// Add to queue
		queueBucket := tx.Bucket(queueBucket)
		return queueBucket.Put([]byte(task.ID), []byte{})
	})
}

func (q *BoltQueue) Pop(ctx context.Context) (*Task, error) {
	var task *Task
	err := q.db.Update(func(tx *bbolt.Tx) error {
		// Get first task from queue
		queueBucket := tx.Bucket(queueBucket)
		cursor := queueBucket.Cursor()
		taskID, _ := cursor.First()
		if taskID == nil {
			return nil
		}

		// Get task data
		tasksBucket := tx.Bucket(tasksBucket)
		taskData := tasksBucket.Get(taskID)
		if taskData == nil {
			return fmt.Errorf("task %s not found", string(taskID))
		}

		// Delete from queue
		if err := queueBucket.Delete(taskID); err != nil {
			return fmt.Errorf("failed to remove task from queue: %v", err)
		}

		// Unmarshal task
		task = &Task{}
		if err := json.Unmarshal(taskData, task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (q *BoltQueue) Get(ctx context.Context, id string) (*Task, error) {
	var task *Task
	err := q.db.View(func(tx *bbolt.Tx) error {
		taskData := tx.Bucket(tasksBucket).Get([]byte(id))
		if taskData == nil {
			return fmt.Errorf("task %s not found", id)
		}

		task = &Task{}
		if err := json.Unmarshal(taskData, task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %v", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (q *BoltQueue) Update(ctx context.Context, id string, status TaskStatus, result *pb.PluginResponse, err string) error {
	return q.db.Update(func(tx *bbolt.Tx) error {
		tasksBucket := tx.Bucket(tasksBucket)
		taskData := tasksBucket.Get([]byte(id))
		if taskData == nil {
			return fmt.Errorf("task %s not found", id)
		}

		task := &Task{}
		if err := json.Unmarshal(taskData, task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %v", err)
		}

		task.Status = status
		task.Result = result
		task.Error = err
		task.UpdatedAt = time.Now()

		taskData, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal updated task: %v", err)
		}

		return tasksBucket.Put([]byte(id), taskData)
	})
}

func (q *BoltQueue) List(ctx context.Context, status *TaskStatus) ([]*Task, error) {
	var tasks []*Task
	err := q.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(tasksBucket).ForEach(func(k, v []byte) error {
			task := &Task{}
			if err := json.Unmarshal(v, task); err != nil {
				return fmt.Errorf("failed to unmarshal task: %v", err)
			}

			if status == nil || task.Status == *status {
				tasks = append(tasks, task)
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (q *BoltQueue) Close() error {
	return q.db.Close()
}
