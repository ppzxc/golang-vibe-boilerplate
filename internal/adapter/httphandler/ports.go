package httphandler

import (
	"context"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

// TodoService is the inbound port defined by this adapter.
type TodoService interface {
	Create(ctx context.Context, title, description string) (*domain.Todo, error)
	FindByID(ctx context.Context, id string) (*domain.Todo, error)
	FindAll(ctx context.Context, cursor string, pageSize int) ([]*domain.Todo, string, int64, error)
	Update(ctx context.Context, id string, title *string, description *string) (*domain.Todo, error)
	Delete(ctx context.Context, id string) error
	Complete(ctx context.Context, id string) (*domain.Todo, error)
}
