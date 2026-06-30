package todo_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/h0ugetsu/todo-api/internal/todo"
)

func TestRepository_Create(t *testing.T) {
	t.Run("created_todo", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		got, err := repo.Create(ctx, todo.Todo{Description: "牛乳を買う"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		if got.ID == 0 {
			t.Errorf("Create() ID = 0, want non-zero")
		}
		if got.Description != "牛乳を買う" {
			t.Errorf("Create() Description = %q, want: %q", got.Description, "牛乳を買う")
		}
		if got.Completed {
			t.Errorf("Create() Completed = true, want false")
		}
		if got.CreatedAt.IsZero() {
			t.Errorf("Create() CreatedAt is zero, want a timestamp")
		}
		if got.UpdatedAt.IsZero() {
			t.Errorf("Create() UpdatedAt is zero, want a timestamp")
		}
	})

	t.Run("must_increment_id", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		first, err := repo.Create(ctx, todo.Todo{Description: "1番目"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}
		second, err := repo.Create(ctx, todo.Todo{Description: "2番目"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		if second.ID != first.ID+1 {
			t.Errorf("Create() second.ID = %d, want %d", second.ID, first.ID+1)
		}
	})
}

func TestRepository_Create_Concurrent(t *testing.T) {
	repo := todo.NewRepository()
	ctx := context.Background()

	const n = 100
	ids := make([]int, n)

	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			got, err := repo.Create(ctx, todo.Todo{Description: "concurrent"})
			if err != nil {
				t.Errorf("Create() returned unexpected error: %v", err)
				return
			}
			ids[i] = got.ID
		})
	}
	wg.Wait()

	seen := make(map[int]bool, n)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID found: %d", id)
		}
		seen[id] = true
	}
}

func TestRepository_FindAll(t *testing.T) {
	t.Run("empty_store", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		got, err := repo.FindAll(ctx)
		if err != nil {
			t.Fatalf("FindAll() returned unexpected error: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("FindAll() is must be empty")
		}
	})

	t.Run("multiple_todos", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		first, err := repo.Create(ctx, todo.Todo{Description: "1番目"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}
		second, err := repo.Create(ctx, todo.Todo{Description: "2番目"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		got, err := repo.FindAll(ctx)
		if err != nil {
			t.Fatalf("FindAll() returned unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("FindAll() len = %d, want 2", len(got))
		}
		if got[0].ID != first.ID {
			t.Errorf("FindAll()[0].ID = %d, want %d", got[0].ID, first.ID)
		}
		if got[1].ID != second.ID {
			t.Errorf("FindAll()[1].ID = %d, want %d", got[1].ID, second.ID)
		}
	})
}

func TestRepository_FindByID(t *testing.T) {
	t.Run("empty_store", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		_, err := repo.FindByID(ctx, 999)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("FindByID() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("found_one", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		created, err := repo.Create(ctx, todo.Todo{Description: "牛乳を買う"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		got, err := repo.FindByID(ctx, created.ID)
		if err != nil {
			t.Fatalf("FindByID() returned unexpected error: %v", err)
		}

		if got.ID != created.ID {
			t.Fatalf("FindByID().ID = %d, want %d", got.ID, created.ID)
		}
	})
}

func TestRepository_Update(t *testing.T) {
	t.Run("updates_description_and_completed", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		created, err := repo.Create(ctx, todo.Todo{Description: "変更前"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		newDescription := "変更後"
		newCompleted := true

		got, err := repo.Update(ctx, created.ID, &newDescription, &newCompleted)
		if err != nil {
			t.Fatalf("Update() returned unexpected error: %v", err)
		}

		if got.Description != newDescription {
			t.Fatalf("Update() Description = %s, want %s", got.Description, newDescription)
		}
		if got.Completed != true {
			t.Fatalf("Update() Completed = %v, want %v", got.Completed, newCompleted)
		}
		if !got.CreatedAt.Equal(created.CreatedAt) {
			t.Errorf("Update() CreatedAt = %v, want unchanged %v", got.CreatedAt, created.CreatedAt)
		}
	})
	t.Run("description_only", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		created, err := repo.Create(ctx, todo.Todo{Description: "変更前"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		newDescription := "変更後"

		got, err := repo.Update(ctx, created.ID, &newDescription, nil)
		if err != nil {
			t.Fatalf("Update() returned unexpected error: %v", err)
		}

		if got.Description != newDescription {
			t.Fatalf("Update() Description = %q, want %q", got.Description, newDescription)
		}
		if got.Completed != created.Completed {
			t.Fatalf("Update() Completed = %v, want unchanged %v", got.Completed, created.Completed)
		}
	})

	t.Run("completed_only", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		created, err := repo.Create(ctx, todo.Todo{Description: "変更前"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		newCompleted := true

		got, err := repo.Update(ctx, created.ID, nil, &newCompleted)
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

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		newDescription := "変更後"

		_, err := repo.Update(ctx, 999, &newDescription, nil)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("Update() err = %v, want ErrNotFound", err)
		}
	})
}

func TestRepository_Delete(t *testing.T) {
	t.Run("delete_todo", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		created, err := repo.Create(ctx, todo.Todo{Description: "牛乳を買う"})
		if err != nil {
			t.Fatalf("Create() returned unexpected error: %v", err)
		}

		err = repo.Delete(ctx, created.ID)
		if err != nil {
			t.Fatalf("Delete() returned unexpected error: %v", err)
		}

		_, err = repo.FindByID(ctx, created.ID)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("FindByID() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		repo := todo.NewRepository()
		ctx := context.Background()

		err := repo.Delete(ctx, 999)
		if !errors.Is(err, todo.ErrNotFound) {
			t.Fatalf("Delete() error = %v, want ErrNotFound", err)
		}
	})
}
