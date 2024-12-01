package queue

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.quinn.io/dataq/dq"
)

type SQLiteQueue struct {
	db *sql.DB
}

// NewSQLiteQueue creates a new SQLite-backed queue
func NewSQLiteQueue(path string) (*SQLiteQueue, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tasks table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			plugin_id TEXT NOT NULL,
			config TEXT NOT NULL,
			data BLOB,
			status TEXT NOT NULL,
			error TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteQueue{db: db}, nil
}

func (q *SQLiteQueue) Push(ctx context.Context, task *Task) error {
	if task.Meta.ID == "" {
		task.Meta.ID = uuid.New().String()
	}
	if task.Meta.CreatedAt.IsZero() {
		task.Meta.CreatedAt = time.Now()
	}
	task.Meta.UpdatedAt = time.Now()

	config, err := json.Marshal(task.Meta.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	_, err = q.db.ExecContext(ctx, `
		INSERT INTO tasks (id, plugin_id, config, data, status, error, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Meta.ID,
		task.Meta.PluginID,
		string(config),
		data,
		string(task.Meta.Status),
		task.Meta.Error,
		task.Meta.CreatedAt,
		task.Meta.UpdatedAt,
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
		SELECT id, plugin_id, config, data, status, error, created_at, updated_at
		FROM tasks
		WHERE status = ?
		ORDER BY created_at ASC
		LIMIT 1
	`, TaskStatusPending)

	var configStr string
	var data []byte
	err = row.Scan(
		&task.Meta.ID,
		&task.Meta.PluginID,
		&configStr,
		&data,
		&task.Meta.Status,
		&task.Meta.Error,
		&task.Meta.CreatedAt,
		&task.Meta.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	err = json.Unmarshal([]byte(configStr), &task.Meta.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Deserialize data using dq package
	if len(data) > 0 {
		task.Data, err = dq.Read(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize data: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM tasks
		WHERE id = ?
	`, task.Meta.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

func (q *SQLiteQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	query := `
		SELECT id, plugin_id, config, data, status, error, created_at, updated_at
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
		var configStr string
		var data []byte
		err = rows.Scan(
			&task.Meta.ID,
			&task.Meta.PluginID,
			&configStr,
			&data,
			&task.Meta.Status,
			&task.Meta.Error,
			&task.Meta.CreatedAt,
			&task.Meta.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		err = json.Unmarshal([]byte(configStr), &task.Meta.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// Deserialize data using dq package
		if len(data) > 0 {
			task.Data, err = dq.Read(bytes.NewReader(data))
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize data: %w", err)
			}
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (q *SQLiteQueue) Update(ctx context.Context, task *Task) error {
	task.Meta.UpdatedAt = time.Now()

	config, err := json.Marshal(task.Meta.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	_, err = q.db.ExecContext(ctx, `
		UPDATE tasks
		SET plugin_id = ?,
			config = ?,
			data = ?,
			status = ?,
			error = ?,
			updated_at = ?
		WHERE id = ?
	`,
		task.Meta.PluginID,
		string(config),
		data,
		string(task.Meta.Status),
		task.Meta.Error,
		task.Meta.UpdatedAt,
		task.Meta.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (q *SQLiteQueue) Close() error {
	return q.db.Close()
}
