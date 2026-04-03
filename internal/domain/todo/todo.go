// Package todo contains the Todo domain model and business logic.
package todo

import (
	"errors"
	"fmt"
	"time"
)

// Sentinel errors for the Todo domain.
var (
	ErrNotFound      = errors.New("todo: not found")
	ErrAlreadyDone   = errors.New("todo: already completed")
	ErrTitleRequired = errors.New("todo: title is required")
)

// Todo is the domain model.
type Todo struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New creates a new Todo with the given title and description.
// Returns ErrTitleRequired if the title is empty.
func New(id, title, description string) (*Todo, error) {
	if title == "" {
		return nil, ErrTitleRequired
	}
	now := time.Now().UTC()
	return &Todo{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Complete marks the Todo as completed.
// Returns ErrAlreadyDone if already completed.
func (t *Todo) Complete() error {
	if t.Completed {
		return ErrAlreadyDone
	}
	t.Completed = true
	t.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateTitle changes the title of the Todo.
// Returns ErrTitleRequired if the new title is empty.
func (t *Todo) UpdateTitle(title string) error {
	if title == "" {
		return ErrTitleRequired
	}
	t.Title = title
	t.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateDescription changes the description of the Todo.
func (t *Todo) UpdateDescription(description string) {
	t.Description = description
	t.UpdatedAt = time.Now().UTC()
}

// String returns a human-readable representation of the Todo.
func (t *Todo) String() string {
	return fmt.Sprintf("Todo{ID: %s, Title: %q, Completed: %v}", t.ID, t.Title, t.Completed)
}
