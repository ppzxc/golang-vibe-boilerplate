package app

import (
	"context"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

// EventBus is the outbound port for publishing domain events.
type EventBus interface {
	Publish(ctx context.Context, events []domain.Event) error
}
