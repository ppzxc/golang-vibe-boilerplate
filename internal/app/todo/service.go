package todo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/app"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/pagination"
)

// Service implements the Todo use case.
type Service struct {
	repo domain.Repository
	bus  app.EventBus
}

// NewService creates a new Service with the given repository and event bus.
func NewService(repo domain.Repository, bus app.EventBus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
	}
}

// Create creates a new Todo.
func (s *Service) Create(ctx context.Context, title, description string) (*domain.Todo, error) {
	id := uuid.New().String()
	todo, err := domain.New(id, title, description)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Create: %w", err)
	}

	if err := s.repo.Save(ctx, todo); err != nil {
		return nil, fmt.Errorf("todo.Service.Create: %w", err)
	}

	if err := s.bus.Publish(ctx, todo.Events); err != nil {
		return nil, fmt.Errorf("todo.Service.Create: %w", err)
	}
	todo.ClearEvents()

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
func (s *Service) FindAll(ctx context.Context, pageToken string, pageSize int) ([]*domain.Todo, string, int64, error) {
	cursor := ""
	if pageToken != "" {
		var err error
		cursor, err = pagination.DecodeToken(pageToken)
		if err != nil {
			return nil, "", 0, fmt.Errorf("todo.Service.FindAll: %w", err)
		}
	}

	if pageSize <= 0 {
		pageSize = pagination.DefaultPageSize
	}

	todos, nextCursor, err := s.repo.FindAll(ctx, cursor, pageSize)
	if err != nil {
		return nil, "", 0, fmt.Errorf("todo.Service.FindAll: %w", err)
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, "", 0, fmt.Errorf("todo.Service.FindAll: %w", err)
	}

	nextToken := ""
	if nextCursor != "" {
		nextToken = pagination.EncodeToken(nextCursor)
	}

	return todos, nextToken, count, nil
}

// Update partially updates a Todo.
func (s *Service) Update(ctx context.Context, id string, title *string, description *string) (*domain.Todo, error) {
	err := s.repo.Update(ctx, id, func(t *domain.Todo) (*domain.Todo, error) {
		if title != nil {
			if err := t.UpdateTitle(*title); err != nil {
				return nil, err
			}
		}
		if description != nil {
			t.UpdateDescription(*description)
		}
		return t, nil
	})
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Update: %w", err)
	}

	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Update: %w", err)
	}

	if err := s.bus.Publish(ctx, updated.Events); err != nil {
		return nil, fmt.Errorf("todo.Service.Update: %w", err)
	}
	updated.ClearEvents()

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
	err := s.repo.Update(ctx, id, func(t *domain.Todo) (*domain.Todo, error) {
		if err := t.Complete(); err != nil {
			return nil, err
		}
		return t, nil
	})
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Complete: %w", err)
	}

	completed, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo.Service.Complete: %w", err)
	}

	if err := s.bus.Publish(ctx, completed.Events); err != nil {
		return nil, fmt.Errorf("todo.Service.Complete: %w", err)
	}
	completed.ClearEvents()

	return completed, nil
}
