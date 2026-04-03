package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/pagination"
)

// TodoHandler handles HTTP requests for the Todo resource.
type TodoHandler struct {
	svc *apptodo.Service
}

// NewTodoHandler creates a new TodoHandler.
func NewTodoHandler(svc *apptodo.Service) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// Create handles POST /todos — creates a new Todo.
// Returns 201 Created with the new resource and a Location header.
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, apptodo.ErrInvalidJSON)
		return
	}

	todo, err := h.svc.Create(r.Context(), apptodo.CreateCommand{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/todos/"+todo.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toTodoResponse(todo))
}

// Get handles GET /todos/{todoId} — retrieves a single Todo.
func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo, err := h.svc.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toTodoResponse(todo))
}

// List handles GET /todos — returns a paginated list of Todos.
// Returns 200 OK with top-level JSON array, Total-Count header, Link header.
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	pageReq, err := pagination.ParseRequest(r)
	if err != nil {
		writeError(w, r, apptodo.ErrInvalidPageSize)
		return
	}

	result, err := h.svc.FindAll(r.Context(), apptodo.ListQuery{
		PageToken: pageReq.PageToken,
		PageSize:  pageReq.PageSize,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	pageResp := &pagination.Response{
		TotalCount: result.TotalCount,
		NextToken:  result.NextToken,
	}
	pageResp.SetHeaders(w, r)

	items := make([]todoResponse, len(result.Items))
	for i, t := range result.Items {
		items[i] = toTodoResponse(t)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

// Update handles PATCH /todos/{todoId} — partially updates a Todo.
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, apptodo.ErrInvalidJSON)
		return
	}

	todo, err := h.svc.Update(r.Context(), id, apptodo.UpdateCommand{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toTodoResponse(todo))
}

// Delete handles DELETE /todos/{todoId} — removes a Todo.
// Returns 204 No Content.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Complete handles POST /todos/{todoId}:complete — marks a Todo as done.
// Colon-syntax custom action per ppzxc RESTful Guidelines (Google AIP-136).
func (h *TodoHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo, err := h.svc.Complete(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toTodoResponse(todo))
}
