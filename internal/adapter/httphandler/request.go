package httphandler

// createTodoRequest is the HTTP request body for creating a Todo.
type createTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// updateTodoRequest is the HTTP request body for updating a Todo.
type updateTodoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}
