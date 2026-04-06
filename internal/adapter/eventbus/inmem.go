package eventbus

import (
	"context"
	"log/slog"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

// InMemoryEventBus is a simple in-memory implementation of the EventBus.
type InMemoryEventBus struct {
	// Handlers could be added here
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events []domain.Event) error {
	for _, e := range events {
		slog.Info("publishing domain event", "type", b.eventName(e), "occurredAt", e.OccurredAt())
		// In a real app, you'd dispatch to registered handlers here
	}
	return nil
}

func (b *InMemoryEventBus) eventName(e domain.Event) string {
	switch e.(type) {
	case domain.Created:
		return "Created"
	case domain.Completed:
		return "Completed"
	case domain.Updated:
		return "Updated"
	default:
		return "Unknown"
	}
}
