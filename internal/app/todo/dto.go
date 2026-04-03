// Package todo contains the Todo use case (application service).
package todo

import (
	"errors"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

// Sentinel errors for the application layer.
var (
	ErrInvalidJSON     = errors.New("invalid JSON body")
	ErrInvalidPageSize = errors.New("invalid pageSize parameter")
)

// CreateCommand holds the input for creating a Todo.
type CreateCommand struct {
	Title       string
	Description string
}

// UpdateCommand holds the input for updating a Todo.
// Only non-nil fields are applied.
type UpdateCommand struct {
	Title       *string
	Description *string
}

// ListQuery holds the input for listing Todos.
type ListQuery struct {
	PageToken string
	PageSize  int
}

// PageResult holds a page of Todos with pagination metadata.
type PageResult struct {
	Items      []*domain.Todo
	NextToken  string
	TotalCount int64
}
