package todo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/h0ugetsu/todo-api/internal/todo"
)

func TestService_CreateTodo(t *testing.T) {
	t.Run("valid_description", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		got, err := service.CreateTodo(ctx, "牛乳を買う")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		if got.Description != "牛乳を買う" {
			t.Errorf("CreateTodo() Description = %q, want %q", got.Description, "牛乳を買う")
		}
		if got.ID == 0 {
			t.Errorf("CreateTodo() ID = 0, want non-zero")
		}
	})

	t.Run("empty_description", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		_, err := service.CreateTodo(ctx, "")
		if !errors.Is(err, todo.ErrInvalidDescription) {
			t.Fatalf("CreateTodo() error = %v, want ErrInvalidDescription", err)
		}
	})
}

func TestService_ListTodos(t *testing.T) {
	t.Run("empty_store", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		got, err := service.ListTodos(ctx)
		if err != nil {
			t.Fatalf("ListTodos() returned unexpected error: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("ListTodos() is must be empty")
		}
	})

	t.Run("multiple_todos", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		first, err := service.CreateTodo(ctx, "1番目")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}
		second, err := service.CreateTodo(ctx, "2番目")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		got, err := service.ListTodos(ctx)
		if err != nil {
			t.Fatalf("ListTodos() returned unexpected error: %v", err)
		}
		if got[0] != first {
			t.Errorf("ListTodos()[0] = %v, got %v", got[0], first)
		}
		if got[1] != second {
			t.Errorf("ListTodos()[1] = %v, got %v", got[1], second)
		}
	})
}

func TestService_GetTodo(t *testing.T) {
	t.Run("found_one", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		created, err := service.CreateTodo(ctx, "牛乳を買う")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		got, err := service.GetTodo(ctx, created.ID)
		if err != nil {
			t.Fatalf("GetTodo() returned unexpected error: %v", err)
		}
		if got.ID != created.ID {
			t.Fatalf("GetTodo() ID = %d, got %d", got.ID, created.ID)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		_, err := service.GetTodo(ctx, 999)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("GetTodo() err = %v, want ErrNotFound", err)
		}
	})
}

func TestService_Update(t *testing.T) {
	t.Run("updates_description_and_completed", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		created, err := service.CreateTodo(ctx, "変更前")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		newDescription := "変更後"
		newCompleted := true

		got, err := service.UpdateTodo(ctx, created.ID, &newDescription, &newCompleted)
		if err != nil {
			t.Fatalf("UpdateTodo() returned unexpected error: %v", err)
		}

		if got.Description != newDescription {
			t.Fatalf("Update() Description = %s, want %s", got.Description, newDescription)
		}
		if got.Completed != true {
			t.Fatalf("Update() Completed = %v, want %v", got.Completed, newCompleted)
		}
	})

	t.Run("description_only", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		created, err := service.CreateTodo(ctx, "変更前")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		newDescription := "変更後"

		got, err := service.UpdateTodo(ctx, created.ID, &newDescription, nil)
		if err != nil {
			t.Fatalf("UpdateTodo() returned unexpected error: %v", err)
		}
		if got.Description != newDescription {
			t.Fatalf("UpdateTodo() Description = %q, want %q", got.Description, newDescription)
		}
		if got.Completed != created.Completed {
			t.Fatalf("UpdateTodo() Completed = %v, want unchanged %v", got.Completed, created.Completed)
		}
	})

	t.Run("completed_only", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		created, err := service.CreateTodo(ctx, "変更前")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		newCompleted := true

		got, err := service.UpdateTodo(ctx, created.ID, nil, &newCompleted)
		if err != nil {
			t.Fatalf("Update() returned unexpected error: %v", err)
		}
		if got.Completed != newCompleted {
			t.Fatalf("Update() Completed = %v, want %v", got.Completed, newCompleted)
		}
		if got.Description != created.Description {
			t.Fatalf("Update() Description = %q, want unchanged %q", got.Description, created.Description)
		}
	})

	t.Run("invalid_description", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		created, err := service.CreateTodo(ctx, "変更前")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		newDescription := ""

		_, err = service.UpdateTodo(ctx, created.ID, &newDescription, nil)
		if !errors.Is(err, todo.ErrInvalidDescription) {
			t.Fatalf("UpdateTodo() err = %v, want ErrInvalidDescription", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		newDescription := "変更後"

		_, err := service.UpdateTodo(ctx, 999, &newDescription, nil)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("UpdateTodo() err = %v, want ErrNotFound", err)
		}
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("delete_todo", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		created, err := service.CreateTodo(ctx, "牛乳を買う")
		if err != nil {
			t.Fatalf("CreateTodo() returned unexpected error: %v", err)
		}

		err = service.DeleteTodo(ctx, created.ID)
		if err != nil {
			t.Fatalf("DeleteTodo() returned unexpected error: %v", err)
		}

		_, err = service.GetTodo(ctx, created.ID)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("GetTodo() err = %v, want ErrNotFound", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		service := todo.NewService(todo.NewRepository())
		ctx := context.Background()

		err := service.DeleteTodo(ctx, 999)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("DeleteTodo() err = %v, want ErrNotFound", err)
		}
	})
}
