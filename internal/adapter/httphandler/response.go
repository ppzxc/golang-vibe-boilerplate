package httphandler

import (
	"time"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

// todoResponse is the HTTP response body for a single Todo.
// All fields are camelCase per ppzxc RESTful Guidelines.
type todoResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func toTodoResponse(t *domain.Todo) todoResponse {
	return todoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
