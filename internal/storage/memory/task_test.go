package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/enes-e-aydogan/stride/internal/modules/task"
	"github.com/enes-e-aydogan/stride/internal/storage/memory"
)

func TestTaskStorage_Create(t *testing.T) {
	t.Run("valid task", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		err := storage.Create(context.Background(), testTask)
		if err != nil {
			t.Fatalf("Unexpected error on create: %v", err)
		}

		_, gotErr := storage.Get(context.Background(), testTask.ID)

		if gotErr != nil {
			t.Fatalf("Expected task to be created, but got error: %v", gotErr)
		}
	})

	t.Run("nil task", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		err := storage.Create(context.Background(), nil)

		if !errors.Is(err, memory.ErrTaskNil) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskNil, err)
		}
	})

	t.Run("existing task", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		err := storage.Create(context.Background(), testTask)
		if err != nil {
			t.Fatalf("Unexpected error on first create: %v", err)
		}

		gotErr := storage.Create(context.Background(), testTask)

		if !errors.Is(gotErr, memory.ErrTaskAlreadyExists) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskAlreadyExists, gotErr)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := storage.Create(ctx, testTask)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
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
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)
		_, err := storage.Get(context.Background(), testTask.ID)
		if err != nil {
			t.Fatalf("Unexpected error on get: %v", err)
		}
	})

	t.Run("empty ID", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		_, err := storage.Get(context.Background(), "")
		if !errors.Is(err, memory.ErrInvalidID) {
			t.Fatalf("Expected error %v, got %v", memory.ErrInvalidID, err)
		}
	})

	t.Run("non-existing ID", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		_, err := storage.Get(context.Background(), "non-existing-id")
		if !errors.Is(err, memory.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskNotFound, err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := storage.Get(ctx, testTask.ID)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

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
		storage := memory.NewTaskStorage()

		testTask, _ := task.NewTask("Test Task")
		testTask1, _ := task.NewTask("Test Task 1")
		testTask2, _ := task.NewTask("Test Task 2")

		storage.Create(context.Background(), testTask)
		storage.Create(context.Background(), testTask1)
		storage.Create(context.Background(), testTask2)

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
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := storage.List(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

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
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Initial Task")

		storage.Create(context.Background(), testTask)
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
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Non-existing Task")

		storage.Create(context.Background(), testTask)
		testTask1, _ := task.NewTask("Another Task")
		testTask1.Title = "Updated Title"

		err := storage.Update(context.Background(), testTask1)
		if !errors.Is(err, memory.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskNotFound, err)
		}
	})

	t.Run("nil task", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		err := storage.Update(context.Background(), nil)
		if !errors.Is(err, memory.ErrTaskNil) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskNil, err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := storage.Update(ctx, testTask)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

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
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

		gotErr := storage.Delete(context.Background(), testTask.ID)

		if gotErr != nil {
			t.Fatalf("Unexpected error on delete: %v", gotErr)
		}

		_, err := storage.Get(context.Background(), testTask.ID)

		if !errors.Is(err, memory.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskNotFound, err)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		storage := memory.NewTaskStorage()

		err := storage.Delete(context.Background(), "")

		if !errors.Is(err, memory.ErrInvalidID) {
			t.Fatalf("Expected error %v, got %v", memory.ErrInvalidID, err)
		}
	})

	t.Run("non-existing-id", func(t *testing.T) {
		storage := memory.NewTaskStorage()

		err := storage.Delete(context.Background(), "non-existing-id")

		if !errors.Is(err, memory.ErrTaskNotFound) {
			t.Fatalf("Expected error %v, got %v", memory.ErrTaskNotFound, err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := storage.Delete(ctx, testTask.ID)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		storage := memory.NewTaskStorage()
		testTask, _ := task.NewTask("Test Task")
		storage.Create(context.Background(), testTask)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(5 * time.Millisecond)

		err := storage.Delete(ctx, testTask.ID)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})
}
