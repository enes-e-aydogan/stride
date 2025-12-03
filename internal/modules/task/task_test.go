package task_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/enes-e-aydogan/stride/internal/modules/task"
)

func TestTask_NewTask_Title(t *testing.T) {
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

func TestTask_WithDescription(t *testing.T) {
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

func TestTask_WithPriority(t *testing.T) {
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

func TestTask_WithDate(t *testing.T) {
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

func TestTask_NewTask_MultipleOptions(t *testing.T) {
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

func TestTask_NewTask_UniqueIDs(t *testing.T) {
	task1, _ := task.NewTask("task1")
	task2, _ := task.NewTask("task2")

	if task1.ID == "" {
		t.Fatalf("expected non-empty ID")
	}
	if task1.ID == task2.ID {
		t.Errorf("expected unique IDs")
	}
}

func TestTask_NewTask_NilOption(t *testing.T) {
	_, err := task.NewTask("Test", nil)
	if err != nil {
		t.Errorf("expected nil option to be ignored, got error: %v", err)
	}
}

func TestTask_NewTask_OptionsAppliedInOrder(t *testing.T) {
	got, _ := task.NewTask("Test",
		task.WithPriority(task.Low),
		task.WithPriority(task.High),
	)
	if got.Priority != task.High {
		t.Errorf("expected last priority to win")
	}
}

func TestTask_TestSetStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      task.Status
		wantErr     bool
		expectedErr error
	}{
		{
			"set to in-progress",
			task.StatusInProgress,

			false,
			nil,
		},
		{
			"set to completed",
			task.StatusCompleted,
			false,
			nil,
		},
		{
			"set to postponed",
			task.StatusPostponed,
			false,
			nil,
		},
		{
			"set to pending",
			task.StatusPending,
			false,
			nil,
		},
		{
			"set to invalid status",
			6,
			true,
			task.ErrInvalidStatus,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := task.NewTask("test task")
			if err != nil {
				t.Fatalf("failed to create task: %v", err)
			}
			gotErr := got.SetStatus(test.status)
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
				if got.Status != test.status {
					t.Errorf("expected status %v, got %v", test.status, got.Status)
				}
				if test.status == task.StatusCompleted {
					if got.CompletedAt == nil {
						t.Errorf("expected CompletedAt to be set, got zero value")
					}
					if got.PostponedAt != nil {
						t.Errorf("expected PostponedAt to be nil, got %v", got.PostponedAt)
					}
				}
				if test.status == task.StatusPostponed {
					if got.CompletedAt != nil {
						t.Errorf("expected CompletedAt to be nil, got %v", got.CompletedAt)
					}
					if got.PostponedAt == nil {
						t.Errorf("expected PostponedAt to be set, got zero value")
					}
				}
				if test.status == task.StatusPending || test.status == task.StatusInProgress {
					if got.CompletedAt != nil {
						t.Errorf("expected CompletedAt to be nil, got %v", got.CompletedAt)
					}
					if got.PostponedAt != nil {
						t.Errorf("expected PostponedAt to be nil, got %v", got.PostponedAt)
					}
				}
				if got.UpdatedAt == nil {
					t.Errorf("expected UpdatedAt to be set, got zero value")
				}
			}
		})
	}
}

func TestTask_CycleStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      task.Status
		newStatus   task.Status
		wantErr     bool
		expectedErr error
	}{
		{
			"pending to in-progress",
			task.StatusPending,
			task.StatusInProgress,

			false,
			nil,
		},
		{
			"in-progress to completed",
			task.StatusInProgress,
			task.StatusCompleted,
			false,
			nil,
		},
		{
			"completed to pending",
			task.StatusCompleted,
			task.StatusPending,
			false,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := task.NewTask("test task")
			if err != nil {
				t.Fatalf("failed to create task: %v", err)
			}
			got.SetStatus(test.status)
			got.CycleStatus()

			if got.Status != test.newStatus {
				t.Errorf("expected status %v, got %v", test.newStatus, got.Status)
			}

			if got.UpdatedAt == nil {
				t.Errorf("expected UpdatedAt to be set, got zero value")
			}
		})
	}
}

func TestTask_SetTitle(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		wantErr     bool
		expectedErr error
	}{
		{
			"valid title",
			"valid title",
			false,
			nil,
		},
		{
			"valid title with spaces",
			"   valid title   ",
			false,
			nil,
		},
		{
			"too long title",
			strings.Repeat("t", task.MaxTitleLength+1),
			true,
			task.ErrTitleTooLong,
		},
		{
			"empty title",
			"",
			true,
			task.ErrTitleEmpty,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := task.NewTask("test task")
			if err != nil {
				t.Fatalf("failed to create task: %v", err)
			}
			gotErr := got.SetTitle(test.title)
			if gotErr != nil {
				if !test.wantErr {
					t.Fatalf("SetTitle() failed unexpectedly: %v", gotErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, gotErr)
				}
				return
			}
			if test.wantErr {
				t.Errorf("SetTitle() succeeded unexpectedly")
			}
			if got.Title != strings.TrimSpace(test.title) {
				t.Errorf("expected title %q, got %q", strings.TrimSpace(test.title), got.Title)
			}

			if got.UpdatedAt == nil {
				t.Errorf("expected UpdatedAt to be set, got zero value")
			}
		})
	}
}

func TestTask_SetDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			"valid description",
			"valid description",
			false,
			nil,
		},
		{
			"valid description with spaces",
			"test task",
			false,
			nil,
		},
		{
			"too long description",
			strings.Repeat("d", task.MaxDescriptionLength+1),
			true,
			task.ErrDescriptionTooLong,
		},
		{
			"empty description",
			"",
			false,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := task.NewTask("test task")
			if err != nil {
				t.Fatalf("failed to create task: %v", err)
			}
			gotErr := got.SetDescription(test.description)
			if gotErr != nil {
				if !test.wantErr {
					t.Fatalf("SetDescription() failed unexpectedly: %v", gotErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, gotErr)
				}
				return
			}
			if test.wantErr {
				t.Errorf("SetDescription() succeeded unexpectedly")
			}
			if got.Description != strings.TrimSpace(test.description) {
				t.Errorf("expected description %q, got %q", strings.TrimSpace(test.description), got.Description)
			}
			if got.UpdatedAt == nil {
				t.Errorf("expected UpdatedAt to be set, got zero value")
			}
		})
	}
}

func TestTask_SetPriority(t *testing.T) {
	tests := []struct {
		name        string
		priority    task.Priority
		wantErr     bool
		expectedErr error
	}{
		{
			"valid priority low",
			task.Low,
			false,
			nil,
		},
		{
			"valid priority high",
			task.High,
			false,
			nil,
		},
		{
			"invalid priority",
			7,
			true,
			task.ErrInvalidPriority,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := task.NewTask("test task")
			if err != nil {
				t.Fatalf("failed to create task: %v", err)
			}
			getErr := got.SetPriority(test.priority)
			if getErr != nil {
				if !test.wantErr {
					t.Fatalf("SetPriority() failed unexpectedly: %v", getErr)
				}
				if !errors.Is(getErr, test.expectedErr) {
					t.Errorf("expected error %v, got %v", test.expectedErr, getErr)
				}
				return
			}
			if test.wantErr {
				t.Errorf("SetPriority() succeeded unexpectedly")
			}
			if got.Priority != test.priority {
				t.Errorf("expected priority %v, got %v", test.priority, got.Priority)
			}
			if got.UpdatedAt == nil {
				t.Errorf("expected UpdatedAt to be set, got zero value")
			}
		})
	}
}

func TestTask_SetDueDate(t *testing.T) {
	t.Run("set due date", func(t *testing.T) {
		task, _ := task.NewTask("Test")
		dueDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
		task.SetDueDate(&dueDate)
		if task.DueDate == nil {
			t.Errorf("expected DueDate to be set, got nil")
		}
		if !task.DueDate.Equal(dueDate) {
			t.Errorf("expected DueDate 2025-12-25, got %v", task.DueDate)
		}
		if task.UpdatedAt == nil {
			t.Errorf("expected UpdatedAt to be set, got zero value")
		}
	})

	t.Run("clear due date", func(t *testing.T) {
		dueDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
		task, _ := task.NewTask("Test", task.WithDueDate(dueDate))
		task.SetDueDate(nil)
		if task.DueDate != nil {
			t.Errorf("expected DueDate to be nil, got %v", task.DueDate)
		}
		if task.UpdatedAt == nil {
			t.Errorf("expected UpdatedAt to be set, got zero value")
		}
	})
}

func TestTask_SetDoDate(t *testing.T) {
	t.Run("set do date", func(t *testing.T) {
		task, _ := task.NewTask("Test")
		doDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
		task.SetDoDate(&doDate)
		if task.DoDate == nil {
			t.Errorf("expected DoDate to be set, got nil")
		}
		if !task.DoDate.Equal(doDate) {
			t.Errorf("expected DoDate 2025-12-25, got %v", task.DoDate)
		}
		if task.UpdatedAt == nil {
			t.Errorf("expected UpdatedAt to be set, got zero value")
		}
	})
	t.Run("clear do date", func(t *testing.T) {
		doDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
		task, _ := task.NewTask("Test", task.WithDoDate(doDate))
		task.SetDoDate(nil)
		if task.DoDate != nil {
			t.Errorf("expected DoDate to be nil, got %v", task.DoDate)
		}
		if task.UpdatedAt == nil {
			t.Errorf("expected UpdatedAt to be set, got zero value")
		}
	})
}
