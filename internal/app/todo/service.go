package todo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/pagination"
)

// Service implements the Todo use case.
// It depends only on the domain Repository port.
type Service struct {
	repo domain.Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new Todo.
func (s *Service) Create(ctx context.Context, cmd CreateCommand) (*domain.Todo, error) {
	id := uuid.New().String()
	todo, err := domain.New(id, cmd.Title, cmd.Description)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Create: %w", err)
	}
	if err := s.repo.Save(ctx, todo); err != nil {
		return nil, fmt.Errorf("todo.Service.Create: %w", err)
	}
	return todo, nil
}

// FindByID retrieves a Todo by its ID.
func (s *Service) FindByID(ctx context.Context, id string) (*domain.Todo, error) {
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.FindByID: %w", err)
	}
	return todo, nil
}

// FindAll retrieves a page of Todos.
func (s *Service) FindAll(ctx context.Context, query ListQuery) (*PageResult, error) {
	cursor := ""
	if query.PageToken != "" {
		var err error
		cursor, err = pagination.DecodeToken(query.PageToken)
		if err != nil {
			return nil, fmt.Errorf("todo.Service.FindAll: %w", err)
		}
	}

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = pagination.DefaultPageSize
	}

	todos, nextCursor, err := s.repo.FindAll(ctx, cursor, pageSize)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.FindAll: %w", err)
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.FindAll: %w", err)
	}

	nextToken := ""
	if nextCursor != "" {
		nextToken = pagination.EncodeToken(nextCursor)
	}

	return &PageResult{
		Items:      todos,
		NextToken:  nextToken,
		TotalCount: count,
	}, nil
}

// Update partially updates a Todo.
func (s *Service) Update(ctx context.Context, id string, cmd UpdateCommand) (*domain.Todo, error) {
	var updated *domain.Todo
	err := s.repo.Update(ctx, id, func(t *domain.Todo) (*domain.Todo, error) {
		if cmd.Title != nil {
			if err := t.UpdateTitle(*cmd.Title); err != nil {
				return nil, err
			}
		}
		if cmd.Description != nil {
			t.UpdateDescription(*cmd.Description)
		}
		return t, nil
	})
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Update: %w", err)
	}
	updated, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Update: %w", err)
	}
	return updated, nil
}

// Delete removes a Todo by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("todo.Service.Delete: %w", err)
	}
	return nil
}

// Complete marks a Todo as completed.
func (s *Service) Complete(ctx context.Context, id string) (*domain.Todo, error) {
	var completed *domain.Todo
	err := s.repo.Update(ctx, id, func(t *domain.Todo) (*domain.Todo, error) {
		if err := t.Complete(); err != nil {
			return nil, err
		}
		return t, nil
	})
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Complete: %w", err)
	}
	completed, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Complete: %w", err)
	}
	return completed, nil
}
