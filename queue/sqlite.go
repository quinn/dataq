package queue

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteQueue struct {
	db *sql.DB
}

// NewSQLiteQueue creates a new SQLite-backed queue
func newSQLiteQueue(path string) (*SQLiteQueue, error) {
	panic("currently broken. Needs to implement PluginRequest fields.")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			error TEXT,
			data_hash TEXT,
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
	_, err := q.db.ExecContext(ctx, `
		INSERT INTO tasks (id, status, error, data_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, task.ID, task.Status, task.Error, task.Hash, task.CreatedAt, task.UpdatedAt)
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
		SELECT id, status, error, data_hash, created_at, updated_at
		FROM tasks
		WHERE status = ?
		LIMIT 1
	`, TaskStatusPending)

	err = row.Scan(
		&task.Status,
		&task.Error,
		&task.Hash,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	// Delete the task from the queue
	_, err = tx.ExecContext(ctx, `
		DELETE FROM tasks
		WHERE id = ?
	`, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

func (q *SQLiteQueue) Update(ctx context.Context, meta *Task) error {
	meta.UpdatedAt = time.Now()

	result, err := q.db.ExecContext(ctx, `
		UPDATE tasks
		SET status = ?, error = ?, data_hash = ?, updated_at = ?
		WHERE id = ?
	`, meta.Status, meta.Error, meta.Hash, meta.UpdatedAt, meta.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task not found: %s", meta.ID())
	}

	return nil
}

func (q *SQLiteQueue) List(ctx context.Context, status TaskStatus) ([]*Task, error) {
	query := `SELECT id, status, error, data_hash, created_at, updated_at FROM tasks`
	args := []interface{}{}

	if status != "" {
		query += ` WHERE status = ?`
		args = append(args, status)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.Status,
			&task.Error,
			&task.Hash,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (q *SQLiteQueue) Close() error {
	return q.db.Close()
}
