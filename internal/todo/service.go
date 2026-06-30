package todo

import "errors"

var ErrInvalidDescription = errors.New("description must not be empty")

type Service interface {
	CreateTodo(description string) (Todo, error)
	ListTodos() ([]Todo, error)
	GetTodo(ID int) (Todo, error)
	UpdateTodo(ID int, description *string, completed *bool) (Todo, error)
	DeleteTodo(ID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTodo(description string) (Todo, error) {
	if description == "" {
		return Todo{}, ErrInvalidDescription
	}

	todo, err := s.repo.Create(Todo{Description: description})
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *service) ListTodos() ([]Todo, error) {
	todos, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *service) GetTodo(ID int) (Todo, error) {
	return s.repo.FindByID(ID)
}

func (s *service) UpdateTodo(ID int, description *string, completed *bool) (Todo, error) {
	if description != nil && *description == "" {
		return Todo{}, ErrInvalidDescription
	}

	todo, err := s.repo.Update(ID, description, completed)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *service) DeleteTodo(ID int) error {
	return s.repo.Delete(ID)
}
