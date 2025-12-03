package task_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/enes-e-aydogan/stride/internal/modules/task"
)

func TestNewTask_Title(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid",
			title:       "Test Task",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "empty",
			title:       "",
			wantErr:     true,
			expectedErr: task.ErrTitleEmpty,
		},
		{
			name:        "too long",
			title:       strings.Repeat("t", task.MaxTitleLength+1),
			wantErr:     true,
			expectedErr: task.ErrTitleTooLong,
		},
		{
			name:        "with spaces",
			title:       "   Test Task   ",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "only spaces",
			title:       "     ",
			wantErr:     true,
			expectedErr: task.ErrTitleEmpty,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := task.NewTask(test.title)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected error %v, got nil", test.expectedErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected no error, got %v", gotErr)
				}
				if got.Title != strings.TrimSpace(test.title) {
					t.Errorf("expected title %q, got %q", strings.TrimSpace(test.title), got.Title)
				}
				if got.Description != "" {
					t.Errorf("expected default description empty, got %q", got.Description)
				}
				if got.Priority != task.None {
					t.Errorf("expected default priority task.None, got %v", got.Priority)
				}
				if got.Status != task.StatusPending {
					t.Errorf("expected default status task.StatusPending, got %v", got.Status)
				}
				if got.CreatedAt.IsZero() {
					t.Errorf("expected CreatedAt to be set, got zero value")
				}
			}
		})
	}
}

func TestWithDescription(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid",
			title:       "Test Task",
			description: "This is a test task.",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "too long",
			title:       "Test Task",
			description: strings.Repeat("d", task.MaxDescriptionLength+1),
			wantErr:     true,
			expectedErr: task.ErrDescriptionTooLong,
		},
		{
			name:        "with spaces",
			title:       "   Test Task   ",
			description: "   This is a test task.   ",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "only spaces",
			title:       "Test Task",
			description: "     ",
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := task.NewTask(test.title, task.WithDescription(test.description))
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected error %v, got nil", test.expectedErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected no error, got %v", gotErr)
				}
				if got.Description != strings.TrimSpace(test.description) {
					t.Errorf("expected description %q, got %q", strings.TrimSpace(test.description), got.Description)
				}
			}
		})
	}
}

func TestWithPriority(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		priority    task.Priority
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid none",
			title:       "Test Task",
			priority:    task.None,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "valid low",
			title:       "Test Task",
			priority:    task.Low,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "valid medium",
			title:       "Test Task",
			priority:    task.Medium,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "valid high",
			title:       "Test Task",
			priority:    task.High,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "valid critical",
			title:       "Test Task",
			priority:    task.Critical,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "valid blocker",
			title:       "Test Task",
			priority:    task.Blocker,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "invalid priority high value",
			title:       "Test Task",
			priority:    6,
			wantErr:     true,
			expectedErr: task.ErrInvalidPriority,
		},
		{
			name:        "invalid priority low value",
			title:       "Test Task",
			priority:    -1,
			wantErr:     true,
			expectedErr: task.ErrInvalidPriority,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := task.NewTask(test.title, task.WithPriority(test.priority))
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected error %v, got nil", test.expectedErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected no error, got %v", gotErr)
				}
				if got.Priority != test.priority {
					t.Errorf("expected priority %v, got %v", test.priority, got.Priority)
				}
			}
		})
	}
}

func TestWithDate(t *testing.T) {
	t.Run("task.WithDueDate", func(t *testing.T) {
		got, gotErr := task.NewTask("Test Task", task.WithDueDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)))
		if gotErr != nil {
			t.Errorf("expected no error, got %v", gotErr)
		}
		if got.DueDate.IsZero() {
			t.Errorf("expected DueDate to be set, got zero value")
		}
		if got.DueDate.Format("2006-01-02") != "2024-01-01" {
			t.Errorf("expected DueDate 2024-01-01, got %v", got.DueDate)
		}
	})

	t.Run("task.WithDoDate", func(t *testing.T) {
		got, gotErr := task.NewTask("Test Task", task.WithDoDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)))
		if gotErr != nil {
			t.Errorf("expected no error, got %v", gotErr)
		}
		if got.DoDate.IsZero() {
			t.Errorf("expected DoDate to be set, got zero value")
		}
		if got.DoDate.Format("2006-01-02") != "2024-01-01" {
			t.Errorf("expected DoDate 2024-01-01, got %v", got.DoDate)
		}
	})
}

func TestNewTask_MultipleOptions(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		priority    task.Priority
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid",
			description: "This is a test task.",
			title:       "Test Task",
			priority:    task.High,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "invalid description",
			title:       "Test Task",
			description: strings.Repeat("d", task.MaxDescriptionLength+1),
			priority:    task.Low,
			wantErr:     true,
			expectedErr: task.ErrDescriptionTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := task.NewTask(test.title, task.WithDescription(test.description), task.WithPriority(test.priority))
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected error %v, got nil", test.expectedErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected no error, got %v", gotErr)
				}
				if got.Description != strings.TrimSpace(test.description) {
					t.Errorf("expected description %q, got %q", strings.TrimSpace(test.description), got.Description)
				}
				if got.Priority != test.priority {
					t.Errorf("expected priority %v, got %v", test.priority, got.Priority)
				}
			}
		})
	}
}

func TestNewTask_UniqueIDs(t *testing.T) {
	task1, _ := task.NewTask("task1")
	task2, _ := task.NewTask("task2")

	if task1.ID == "" {
		t.Fatalf("expected non-empty ID")
	}
	if task1.ID == task2.ID {
		t.Errorf("expected unique IDs")
	}
}

func TestNewTask_NilOption(t *testing.T) {
	_, err := task.NewTask("Test", nil)
	if err != nil {
		t.Errorf("expected nil option to be ignored, got error: %v", err)
	}
}

func TestNewTask_OptionsAppliedInOrder(t *testing.T) {
	got, _ := task.NewTask("Test",
		task.WithPriority(task.Low),
		task.WithPriority(task.High),
	)
	if got.Priority != task.High {
		t.Errorf("expected last priority to win")
	}
}
