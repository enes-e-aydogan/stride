package sqlite_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/enes-e-aydogan/stride/internal/modules/task"
	"github.com/enes-e-aydogan/stride/internal/storage/sqlite"
)

func testSetup(t *testing.T, title string) (*sqlite.TaskStorage, *task.Task, error) {
	t.Helper()
	db, connectionErr := sqlite.NewConnection(context.Background(), ":memory:")
	if connectionErr != nil {
		return nil, nil, connectionErr
	}

	storage, storageErr := sqlite.NewTaskStorage(context.Background(), db)
	if storageErr != nil {
		return nil, nil, storageErr
	}

	testTask, taskErr := task.NewTask(title)
	if taskErr != nil {
		return nil, nil, taskErr
	}

	return storage, testTask, nil
}

func TestTaskStorage_Create(t *testing.T) {
	t.Run("valid task", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("Failed to create task in DB: %v", createErr)
		}

		_, gotErr := storage.Get(context.Background(), testTask.ID)
		if gotErr != nil {
			t.Fatalf("Failed to get task: %v", gotErr)
		}
	})

	t.Run("nil task", func(t *testing.T) {
		storage, _, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		err := storage.Create(context.Background(), nil)
		if !errors.Is(err, sqlite.ErrTaskNil) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskNil, err)
		}
	})

	t.Run("existing task", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		err := storage.Create(context.Background(), testTask)
		if err != nil {
			t.Fatalf("Unexpected error on first create: %v", err)
		}

		gotErr := storage.Create(context.Background(), testTask)

		if !errors.Is(gotErr, sqlite.ErrTaskAlreadyExists) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskAlreadyExists, gotErr)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := storage.Create(ctx, testTask)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(5 * time.Millisecond)

		err := storage.Create(ctx, testTask)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})
}

func TestTaskStorage_Get(t *testing.T) {
	t.Run("valid ID", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}
		_, err := storage.Get(context.Background(), testTask.ID)
		if err != nil {
			t.Fatalf("Unexpected error on get: %v", err)
		}
	})

	t.Run("empty ID", func(t *testing.T) {
		storage, _, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		_, err := storage.Get(context.Background(), "")
		if !errors.Is(err, sqlite.ErrInvalidID) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrInvalidID, err)
		}
	})

	t.Run("non-existing ID", func(t *testing.T) {
		storage, _, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		_, err := storage.Get(context.Background(), "non-existing-id")
		if !errors.Is(err, sqlite.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskNotFound, err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := storage.Get(ctx, testTask.ID)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(5 * time.Millisecond)
		_, err := storage.Get(ctx, testTask.ID)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})
}

func TestTaskStorage_List(t *testing.T) {
	t.Run("list tasks", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		testTask1, _ := task.NewTask("Test Task 1")
		testTask2, _ := task.NewTask("Test Task 2")

		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}
		createErr = storage.Create(context.Background(), testTask1)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}
		createErr = storage.Create(context.Background(), testTask2)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		tasks, err := storage.List(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error on list: %v", err)
		}

		if len(tasks) != 3 {
			t.Fatalf("Expected 3 tasks, got %d", len(tasks))
		}

		for _, expected := range []*task.Task{testTask, testTask1, testTask2} {
			found := false
			for _, got := range tasks {
				if got.ID == expected.ID {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Task %s not found in list", expected.ID)
			}
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := storage.List(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(5 * time.Millisecond)

		_, err := storage.List(ctx)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})
}

func TestTaskStorage_Update(t *testing.T) {
	t.Run("valid update", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}
		testTask.Title = "Updated Task"

		err := storage.Update(context.Background(), testTask)
		if err != nil {
			t.Fatalf("Unexpected error on update: %v", err)
		}

		got, err := storage.Get(context.Background(), testTask.ID)
		if err != nil {
			t.Fatalf("Unexpected error on get after update: %v", err)
		}

		if got.Title != "Updated Task" {
			t.Fatalf("Expected title 'Updated Task', got '%s'", got.Title)
		}
	})

	t.Run("non-existing task", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}
		testTask1, _ := task.NewTask("Another Task")
		testTask1.Title = "Updated Title"

		err := storage.Update(context.Background(), testTask1)
		if !errors.Is(err, sqlite.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskNotFound, err)
		}
	})

	t.Run("nil task", func(t *testing.T) {
		storage, _, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		err := storage.Update(context.Background(), nil)
		if !errors.Is(err, sqlite.ErrTaskNil) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskNil, err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := storage.Update(ctx, testTask)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(5 * time.Millisecond)

		err := storage.Update(ctx, testTask)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})
}

func TestTaskStorage_Delete(t *testing.T) {
	t.Run("valid delete", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		gotErr := storage.Delete(context.Background(), testTask.ID)

		if gotErr != nil {
			t.Fatalf("Unexpected error on delete: %v", gotErr)
		}

		_, err := storage.Get(context.Background(), testTask.ID)

		if !errors.Is(err, sqlite.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskNotFound, err)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		storage, _, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		err := storage.Delete(context.Background(), "")

		if !errors.Is(err, sqlite.ErrInvalidID) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrInvalidID, err)
		}
	})

	t.Run("non-existing-id", func(t *testing.T) {
		storage, _, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}

		err := storage.Delete(context.Background(), "non-existing-id")

		if !errors.Is(err, sqlite.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", sqlite.ErrTaskNotFound, err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := storage.Delete(ctx, testTask.ID)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage, testTask, setupErr := testSetup(t, "Test Task")
		if setupErr != nil {
			t.Fatalf("test setup failed: %v", setupErr)
		}
		createErr := storage.Create(context.Background(), testTask)
		if createErr != nil {
			t.Fatalf("failed to create task: %v", createErr)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(5 * time.Millisecond)

		err := storage.Delete(ctx, testTask.ID)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})
}
