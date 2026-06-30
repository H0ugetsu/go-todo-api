package todo

import (
	"context"
	"errors"
)

var ErrInvalidDescription = errors.New("description must not be empty")

type Service interface {
	CreateTodo(ctx context.Context, description string) (Todo, error)
	ListTodos(ctx context.Context) ([]Todo, error)
	GetTodo(ctx context.Context, ID int) (Todo, error)
	UpdateTodo(ctx context.Context, ID int, description *string, completed *bool) (Todo, error)
	DeleteTodo(ctx context.Context, ID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTodo(ctx context.Context, description string) (Todo, error) {
	if description == "" {
		return Todo{}, ErrInvalidDescription
	}

	todo, err := s.repo.Create(ctx, Todo{Description: description})
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *service) ListTodos(ctx context.Context) ([]Todo, error) {
	todos, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *service) GetTodo(ctx context.Context, ID int) (Todo, error) {
	return s.repo.FindByID(ctx, ID)
}

func (s *service) UpdateTodo(ctx context.Context, ID int, description *string, completed *bool) (Todo, error) {
	if description != nil && *description == "" {
		return Todo{}, ErrInvalidDescription
	}

	todo, err := s.repo.Update(ctx, ID, description, completed)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *service) DeleteTodo(ctx context.Context, ID int) error {
	return s.repo.Delete(ctx, ID)
}
