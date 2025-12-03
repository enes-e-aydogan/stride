// Package task contains the domain logic for tasks.
package task

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Constants for maximum lengths.
const (
	MaxTitleLength       = 255
	MaxDescriptionLength = 4000
)

// Error definitions for task validation.
var (
	ErrTitleEmpty         = errors.New("title cannot be empty")
	ErrTitleTooLong       = errors.New("title exceeds maximum length")
	ErrDescriptionTooLong = errors.New("description exceeds maximum length")
	ErrInvalidPriority    = errors.New("invalid priority value")
	ErrInvalidStatus      = errors.New("invalid status value")
	ErrCreatedAtNotSet    = errors.New("created at timestamp must be set")
)

// Status represents the status of a task.
type Status int

// Constants for task statuses.
const (
	StatusPending Status = iota
	StatusInProgress
	StatusCompleted
	StatusPostponed
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "Pending"
	case StatusInProgress:
		return "In Progress"
	case StatusCompleted:
		return "Completed"
	case StatusPostponed:
		return "Postponed"
	default:
		return "Unknown"
	}
}

// Priority represents the priority level of a task.
type Priority int

// Constants for task priority levels.
const (
	None Priority = iota
	Low
	Medium
	High
	Critical
	Blocker
)

func (p Priority) String() string {
	switch p {
	case None:
		return "None"
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	case Blocker:
		return "Blocker"
	default:
		return "Unknown"
	}
}

// Task represents a task with its associated data.
type Task struct {
	ID          string
	Title       string
	Description string
	Priority    Priority
	Status      Status
	DoDate      *time.Time
	DueDate     *time.Time
	CompletedAt *time.Time
	PostponedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func validateTitle(title string) error {
	if len(title) == 0 {
		return ErrTitleEmpty
	}
	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > MaxDescriptionLength {
		return ErrDescriptionTooLong
	}
	return nil
}

func validatePriority(p Priority) error {
	if p < None || p > Blocker {
		return ErrInvalidPriority
	}
	return nil
}

func validateStatus(s Status) error {
	if s < StatusPending || s > StatusPostponed {
		return ErrInvalidStatus
	}
	return nil
}

// Option defines a functional option for configuring a Task.
type Option func(*Task) error

// WithDescription sets the Description of the task.
func WithDescription(description string) Option {
	return func(t *Task) error {
		err := validateDescription(description)
		if err != nil {
			return err
		}
		t.Description = strings.TrimSpace(description)
		return nil
	}
}

// WithPriority sets the Priority of the task.
func WithPriority(priority Priority) Option {
	return func(t *Task) error {
		err := validatePriority(priority)
		if err != nil {
			return err
		}
		t.Priority = priority
		return nil
	}
}

// WithDueDate sets the DueDate of the task.
func WithDueDate(dueDate time.Time) Option {
	return func(t *Task) error {
		t.DueDate = &dueDate
		return nil
	}
}

// WithDoDate sets the DoDate of the task.
func WithDoDate(doDate time.Time) Option {
	return func(t *Task) error {
		t.DoDate = &doDate
		return nil
	}
}

// NewTask creates a new Task with the given title and options.
func NewTask(title string, options ...Option) (*Task, error) {
	title = strings.TrimSpace(title)
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	task := &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: "",
		Priority:    None,
		Status:      StatusPending,
		CreatedAt:   time.Now().UTC(),
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(task); err != nil {
			return nil, err
		}
	}

	return task, nil
}
