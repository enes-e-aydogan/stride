package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/enes-e-aydogan/stride/internal/modules/task"
)

// Error definitions for task storage operations.
var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskNil           = errors.New("task cannot be nil")
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrInvalidID         = errors.New("invalid task ID")
)

// TaskStorage is a SQLite implementation of the Storage interface.
type TaskStorage struct {
	db *sql.DB
}

// NewTaskStorage creates a new instance of TaskStorage and returns its pointer.
func NewTaskStorage(ctx context.Context, db *sql.DB) (*TaskStorage, error) {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	priority INTEGER NOT NULL,
	status INTEGER NOT NULL,
	do_date TIMESTAMP,
	due_date TIMESTAMP,
	completed_at TIMESTAMP,
	postponed_at TIMESTAMP,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP)
	`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &TaskStorage{db: db}, nil
}

func toNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func fromNullTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// Create adds a new task to the storage.
func (ts *TaskStorage) Create(ctx context.Context, task *task.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if task == nil {
		return ErrTaskNil
	}

	_, err := ts.db.ExecContext(ctx, `
         INSERT INTO tasks (
             id, title, description, priority, status,
             do_date, due_date, completed_at, postponed_at,
             created_at, updated_at
         ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
         `,
		task.ID,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		toNullTime(task.DoDate),
		toNullTime(task.DueDate),
		toNullTime(task.CompletedAt),
		toNullTime(task.PostponedAt),
		task.CreatedAt,
		toNullTime(task.UpdatedAt),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrTaskAlreadyExists
		}
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

// Get returns a task by its ID.
func (ts *TaskStorage) Get(ctx context.Context, id string) (*task.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if id == "" {
		return nil, ErrInvalidID
	}

	res := task.Task{}
	var doDate, dueDate, completedAt, postponedAt, updatedAt sql.NullTime
	err := ts.db.QueryRowContext(ctx, `SELECT id, title, description,
		priority, status, do_date, due_date,
		completed_at, postponed_at, created_at, updated_at from tasks WHERE id = ?`, id).
		Scan(&res.ID, &res.Title, &res.Description, &res.Priority, &res.Status,
			&doDate, &dueDate, &completedAt, &postponedAt, &res.CreatedAt, &updatedAt)
	res.DoDate = fromNullTime(doDate)
	res.DueDate = fromNullTime(dueDate)
	res.CompletedAt = fromNullTime(completedAt)
	res.PostponedAt = fromNullTime(postponedAt)
	res.UpdatedAt = fromNullTime(updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &res, nil
}

// List returns all tasks in the storage.
func (ts *TaskStorage) List(ctx context.Context) ([]*task.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	res := []*task.Task{}

	rows, err := ts.db.QueryContext(ctx, `SELECT id, title, description,
		priority, status, do_date, due_date,
		completed_at, postponed_at, created_at, updated_at from tasks ORDER BY created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query task: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := &task.Task{}
		var doDate, dueDate, completedAt, postponedAt, updatedAt sql.NullTime
		rowsErr := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.Status,
			&doDate, &dueDate, &completedAt, &postponedAt, &t.CreatedAt, &updatedAt)
		if rowsErr != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", rowsErr)
		}

		t.DoDate = fromNullTime(doDate)
		t.DueDate = fromNullTime(dueDate)
		t.CompletedAt = fromNullTime(completedAt)
		t.PostponedAt = fromNullTime(postponedAt)
		t.UpdatedAt = fromNullTime(updatedAt)
		res = append(res, t)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rowsErr)
	}
	return res, nil
}

// Update modifies an existing task in the storage.
func (ts *TaskStorage) Update(ctx context.Context, task *task.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if task == nil {
		return ErrTaskNil
	}

	result, err := ts.db.ExecContext(ctx, `
             UPDATE tasks SET
                 title = ?, description = ?, priority = ?, status = ?,
                 do_date = ?, due_date = ?,
                 completed_at = ?, postponed_at = ?, updated_at = ?
             WHERE id = ?
         `, task.Title, task.Description, task.Priority, task.Status,
		toNullTime(task.DoDate), toNullTime(task.DueDate),
		toNullTime(task.CompletedAt), toNullTime(task.PostponedAt),
		toNullTime(task.UpdatedAt), task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// Delete removes a task from the storage by its ID.
func (ts *TaskStorage) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if id == "" {
		return ErrInvalidID
	}

	result, err := ts.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}
