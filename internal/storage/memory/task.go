// Package memory implements an in-memory storage backend.
package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/enes-e-aydogan/stride/internal/modules/task"
)

// Error definitions for task storage operations.
var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskNil           = errors.New("task cannot be nil")
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrInvalidID         = errors.New("invalid task ID")
)

// TaskStorage is an in-memory implementation of the Storage interface.
type TaskStorage struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

// NewTaskStorage creates a new instance of TaskStorage and returns its pointer.
func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		mu:    sync.RWMutex{},
		tasks: make(map[string]*task.Task),
	}
}

// Create adds a new task to the storage.
func (ts *TaskStorage) Create(ctx context.Context, task *task.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if task == nil {
		return ErrTaskNil
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[task.ID]; exists {
		return ErrTaskAlreadyExists
	}

	ts.tasks[task.ID] = task

	return nil
}

// Get retrieves a task by its ID.
func (ts *TaskStorage) Get(ctx context.Context, id string) (*task.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if id == "" {
		return nil, ErrInvalidID
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()
	t, exists := ts.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return t.Copy(), nil
}

// List returns all tasks in the storage.
func (ts *TaskStorage) List(ctx context.Context) ([]*task.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	results := make([]*task.Task, 0, len(ts.tasks))

	for _, t := range ts.tasks {
		results = append(results, t.Copy())
	}

	return results, nil
}

// Update modifies an existing task in the storage.
func (ts *TaskStorage) Update(ctx context.Context, task *task.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if task == nil {
		return ErrTaskNil
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	ts.tasks[task.ID] = task

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

	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(ts.tasks, id)

	return nil
}
