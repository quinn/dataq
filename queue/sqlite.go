package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type SQLiteQueue struct {
	db *sql.DB
}

// NewSQLiteQueue creates a new SQLite-backed queue
func NewSQLiteQueue(opts ...Option) (*SQLiteQueue, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Path == "" {
		return nil, fmt.Errorf("path is required")
	}

	db, err := sql.Open("sqlite3", options.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tasks table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			plugin_id TEXT NOT NULL,
			config TEXT NOT NULL,
			data TEXT NOT NULL,
			status TEXT NOT NULL,
			result TEXT,
			error TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tasks table: %w", err)
	}

	return &SQLiteQueue{db: db}, nil
}

func (q *SQLiteQueue) Push(ctx context.Context, task *Task) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	task.UpdatedAt = time.Now()

	config, err := json.Marshal(task.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	data, err := protojson.Marshal(task.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var result []byte
	if task.Result != nil {
		result, err = protojson.Marshal(task.Result)
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}
	}

	_, err = q.db.ExecContext(ctx, `
		INSERT INTO tasks (id, plugin_id, config, data, status, result, error, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.ID,
		task.PluginID,
		string(config),
		string(data),
		string(task.Status),
		string(result),
		task.Error,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

func (q *SQLiteQueue) Pop(ctx context.Context) (*Task, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var task Task
	row := tx.QueryRowContext(ctx, `
		SELECT id, plugin_id, config, data, status, result, error, created_at, updated_at
		FROM tasks
		WHERE status = ?
		ORDER BY created_at ASC
		LIMIT 1
	`, TaskStatusPending)

	var configStr, dataStr, resultStr string
	err = row.Scan(
		&task.ID,
		&task.PluginID,
		&configStr,
		&dataStr,
		&task.Status,
		&resultStr,
		&task.Error,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	err = json.Unmarshal([]byte(configStr), &task.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	task.Data = &pb.DataItem{}
	err = protojson.Unmarshal([]byte(dataStr), task.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	if resultStr != "" {
		task.Result = &pb.PluginResponse{}
		err = protojson.Unmarshal([]byte(resultStr), task.Result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM tasks
		WHERE id = ?
	`, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

func (q *SQLiteQueue) Update(ctx context.Context, task *Task) error {
	task.UpdatedAt = time.Now()

	config, err := json.Marshal(task.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	data, err := protojson.Marshal(task.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var result []byte
	if task.Result != nil {
		result, err = protojson.Marshal(task.Result)
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}
	}

	_, err = q.db.ExecContext(ctx, `
		UPDATE tasks
		SET plugin_id = ?,
			config = ?,
			data = ?,
			status = ?,
			result = ?,
			error = ?,
			updated_at = ?
		WHERE id = ?
	`,
		task.PluginID,
		string(config),
		string(data),
		string(task.Status),
		string(result),
		task.Error,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (q *SQLiteQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	query := `
		SELECT id, plugin_id, config, data, status, result, error, created_at, updated_at
		FROM tasks
	`
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, string(status))
	}
	query += " ORDER BY created_at ASC"

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var configStr, dataStr, resultStr string
		err = rows.Scan(
			&task.ID,
			&task.PluginID,
			&configStr,
			&dataStr,
			&task.Status,
			&resultStr,
			&task.Error,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		err = json.Unmarshal([]byte(configStr), &task.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		task.Data = &pb.DataItem{}
		err = protojson.Unmarshal([]byte(dataStr), task.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		if resultStr != "" {
			task.Result = &pb.PluginResponse{}
			err = protojson.Unmarshal([]byte(resultStr), task.Result)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal result: %w", err)
			}
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (q *SQLiteQueue) Close() error {
	return q.db.Close()
}
