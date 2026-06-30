package todo

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("todo not found")

type Repository interface {
	Create(t Todo) (Todo, error)
	FindAll() ([]Todo, error)
	FindByID(ID int) (Todo, error)
	Update(ID int, description *string, completed *bool) (Todo, error)
	Delete(ID int) error
}

type repository struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int
}

func NewRepository() Repository {
	return &repository{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (r *repository) Create(t Todo) (Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo := Todo{
		ID:          r.nextID,
		Description: t.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.todos[todo.ID] = todo
	r.nextID++

	return todo, nil
}

func (r *repository) FindAll() ([]Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]Todo, 0, len(r.todos))
	for _, t := range r.todos {
		todos = append(todos, t)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

func (r *repository) FindByID(ID int) (Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.todos[ID]
	if !ok {
		return Todo{}, ErrNotFound
	}

	return todo, nil
}

func (r *repository) Update(ID int, description *string, completed *bool) (Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, ok := r.todos[ID]
	if !ok {
		return Todo{}, ErrNotFound
	}

	if description != nil {
		todo.Description = *description
	}
	if completed != nil {
		todo.Completed = *completed
	}
	todo.UpdatedAt = time.Now()

	r.todos[todo.ID] = todo

	return todo, nil
}

func (r *repository) Delete(ID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.todos[ID]
	if !ok {
		return ErrNotFound
	}

	delete(r.todos, ID)

	return nil
}
