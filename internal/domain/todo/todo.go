package todo

import (
	"errors"
	"time"
)

// Event is the interface for all domain events.
type Event interface {
	OccurredAt() time.Time
}

// Completed is a domain event.
type Completed struct {
	ID   string
	When time.Time
}

func (e Completed) OccurredAt() time.Time { return e.When }

// Created is a domain event.
type Created struct {
	ID   string
	When time.Time
}

func (e Created) OccurredAt() time.Time { return e.When }

// Updated is a domain event.
type Updated struct {
	ID   string
	When time.Time
}

func (e Updated) OccurredAt() time.Time { return e.When }

// Sentinel errors.
var (
	ErrNotFound      = errors.New("todo: not found")
	ErrAlreadyDone   = errors.New("todo: already completed")
	ErrTitleRequired = errors.New("todo: title is required")
)

// Todo is the domain aggregate.
type Todo struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Events      []Event // Domain events
}

func New(id, title, description string) (*Todo, error) {
	if title == "" {
		return nil, ErrTitleRequired
	}
	now := time.Now().UTC()
	t := &Todo{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	t.Events = append(t.Events, Created{ID: id, When: now})
	return t, nil
}

func (t *Todo) Complete() error {
	if t.Completed {
		return ErrAlreadyDone
	}
	t.Completed = true
	t.UpdatedAt = time.Now().UTC()
	t.Events = append(t.Events, Completed{ID: t.ID, When: t.UpdatedAt})
	return nil
}

func (t *Todo) UpdateTitle(title string) error {
	if title == "" {
		return ErrTitleRequired
	}
	t.Title = title
	t.UpdatedAt = time.Now().UTC()
	t.Events = append(t.Events, Updated{ID: t.ID, When: t.UpdatedAt})
	return nil
}

func (t *Todo) UpdateDescription(description string) {
	t.Description = description
	t.UpdatedAt = time.Now().UTC()
	t.Events = append(t.Events, Updated{ID: t.ID, When: t.UpdatedAt})
}

// ClearEvents empties the event list.
func (t *Todo) ClearEvents() {
	t.Events = nil
}

// String returns a human-readable representation of the Todo.
func (t *Todo) String() string {
	return "Todo{ID: " + t.ID + ", Title: " + t.Title + "}"
}
