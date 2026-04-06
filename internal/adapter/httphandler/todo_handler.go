package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/pagination"
)

type TodoHandler struct {
	service TodoService
}

func NewTodoHandler(service TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request, params ListTodosParams) {
	ctx := r.Context()

	pageSize := pagination.DefaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	pageToken := ""
	if params.PageToken != nil {
		pageToken = *params.PageToken
	}

	todos, nextToken, totalCount, err := h.service.FindAll(ctx, pageToken, pageSize)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	resp := make([]Todo, len(todos))
	for i, t := range todos {
		resp[i] = mapDomainToTodo(t)
	}

	pageResp := &pagination.Response{
		TotalCount: totalCount,
		NextToken:  nextToken,
	}
	pageResp.SetHeaders(w, r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}

	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}

	todo, err := h.service.Create(ctx, req.Title, desc)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	resp := mapDomainToTodo(todo)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/todos/"+todo.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request, todoID string) {
	ctx := r.Context()
	todo, err := h.service.FindByID(ctx, todoID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	resp := mapDomainToTodo(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request, todoID string) {
	ctx := r.Context()
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}

	todo, err := h.service.Update(ctx, todoID, req.Title, req.Description)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	resp := mapDomainToTodo(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request, todoID string) {
	ctx := r.Context()
	if err := h.service.Delete(ctx, todoID); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request, todoID string) {
	ctx := r.Context()
	todo, err := h.service.Complete(ctx, todoID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	resp := mapDomainToTodo(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func mapDomainToTodo(t *domain.Todo) Todo {
	var desc *string
	if t.Description != "" {
		desc = &t.Description
	}

	u, _ := uuid.Parse(t.ID)

	return Todo{
		Id:          openapi_types.UUID(u),
		Title:       t.Title,
		Description: desc,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
