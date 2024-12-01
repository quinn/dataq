package queue

import (
	"bytes"
	"context"
	"database/sql"
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
func newSQLiteQueue(path string) (*SQLiteQueue, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tasks table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			plugin_id TEXT NOT NULL,
			status TEXT NOT NULL,
			data BLOB,
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

	// Serialize data using dq package
	var buf bytes.Buffer
	if err := dq.WriteDataItem(&buf, task.Data); err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	data := buf.Bytes()

	_, err := q.db.ExecContext(ctx, `
		INSERT INTO tasks (id, plugin_id, status, created_at, updated_at, data)
		VALUES (?, ?, ?, ?, ?, ?)
	`, task.Meta.ID, task.Data.Meta.PluginId, task.Meta.Status, task.Meta.CreatedAt, task.Meta.UpdatedAt, data)
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
		SELECT id, plugin_id, status, data, created_at, updated_at
		FROM tasks
		WHERE status = ?
		ORDER BY created_at ASC
		LIMIT 1
	`, TaskStatusPending)

	var data []byte
	err = row.Scan(
		&task.Meta.ID,
		&task.Data.Meta.PluginId,
		&task.Meta.Status,
		&data,
		&task.Meta.CreatedAt,
		&task.Meta.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
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

func (q *SQLiteQueue) List(ctx context.Context, status TaskStatus) ([]*TaskMetadata, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT id, plugin_id, status, created_at, updated_at
		FROM tasks
		WHERE status = ? OR ? = ''
	`, status, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*TaskMetadata
	for rows.Next() {
		var task TaskMetadata
		var pluginID string
		err := rows.Scan(&task.ID, &pluginID, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (q *SQLiteQueue) Update(ctx context.Context, meta *TaskMetadata) error {
	meta.UpdatedAt = time.Now()

	_, err := q.db.ExecContext(ctx, `
		UPDATE tasks
		SET status = ?,
			updated_at = ?
		WHERE id = ?
	`,
		meta.Status,
		meta.UpdatedAt,
		meta.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (q *SQLiteQueue) Close() error {
	return q.db.Close()
}
