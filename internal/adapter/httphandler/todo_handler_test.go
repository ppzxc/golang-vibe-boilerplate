package httphandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/httphandler"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTodoRepo is an in-memory implementation of domain.Repository for handler tests.
type mockTodoRepo struct {
	todos map[string]*domain.Todo
}

func newMockTodoRepo() *mockTodoRepo {
	return &mockTodoRepo{todos: make(map[string]*domain.Todo)}
}

func (m *mockTodoRepo) Save(_ context.Context, t *domain.Todo) error {
	m.todos[t.ID] = t
	return nil
}

func (m *mockTodoRepo) FindByID(_ context.Context, id string) (*domain.Todo, error) {
	t, ok := m.todos[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (m *mockTodoRepo) FindAll(_ context.Context, _ string, _ int) ([]*domain.Todo, string, error) {
	var result []*domain.Todo
	for _, t := range m.todos {
		cp := *t
		result = append(result, &cp)
	}
	return result, "", nil
}

func (m *mockTodoRepo) Update(_ context.Context, id string, fn func(*domain.Todo) (*domain.Todo, error)) error {
	t, ok := m.todos[id]
	if !ok {
		return domain.ErrNotFound
	}
	updated, err := fn(t)
	if err != nil {
		return err
	}
	m.todos[id] = updated
	return nil
}

func (m *mockTodoRepo) Delete(_ context.Context, id string) error {
	if _, ok := m.todos[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.todos, id)
	return nil
}

func (m *mockTodoRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.todos)), nil
}

// setupRouter creates a test router with the handler backed by the given repo.
func setupRouter(repo domain.Repository) http.Handler {
	svc := apptodo.NewService(repo)
	return httphandler.NewRouter(svc)
}

func TestHandler_Create(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	body := `{"title":"Buy groceries","description":"milk, eggs"}`
	r := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, w.Header().Get("Location"))
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "Buy groceries", resp["title"])
	assert.NotEmpty(t, resp["id"])
}

func TestHandler_Create_EmptyTitle(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	body := `{"title":""}`
	r := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/problem+json")
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	body := `not valid json`
	r := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/problem+json")
}

func TestHandler_Get(t *testing.T) {
	repo := newMockTodoRepo()
	td, _ := domain.New("id-1", "Test Todo", "some desc")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/todos/id-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "id-1", resp["id"])
	assert.Equal(t, "Test Todo", resp["title"])
}

func TestHandler_Get_NotFound(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/todos/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/problem+json")
}

func TestHandler_List(t *testing.T) {
	repo := newMockTodoRepo()
	td, _ := domain.New("id-1", "Test Todo", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Total-Count"))

	var items []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&items))
	assert.Len(t, items, 1)
}

func TestHandler_List_Empty(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var items []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&items))
	assert.Len(t, items, 0)
}

func TestHandler_Update(t *testing.T) {
	repo := newMockTodoRepo()
	td, _ := domain.New("id-1", "Old Title", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	body := `{"title":"New Title"}`
	r := httptest.NewRequest(http.MethodPatch, "/todos/id-1", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "New Title", resp["title"])
}

func TestHandler_Update_NotFound(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	body := `{"title":"New Title"}`
	r := httptest.NewRequest(http.MethodPatch, "/todos/nonexistent", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Delete(t *testing.T) {
	repo := newMockTodoRepo()
	td, _ := domain.New("id-1", "To delete", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodDelete, "/todos/id-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodDelete, "/todos/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Complete(t *testing.T) {
	repo := newMockTodoRepo()
	td, _ := domain.New("id-1", "Test", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodPost, "/todos/id-1:complete", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, true, resp["completed"])
}

func TestHandler_Complete_AlreadyDone(t *testing.T) {
	repo := newMockTodoRepo()
	td, _ := domain.New("id-1", "Test", "")
	_ = td.Complete()
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodPost, "/todos/id-1:complete", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/problem+json")
}

func TestHandler_Complete_NotFound(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodPost, "/todos/nonexistent:complete", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMiddleware_RequestID(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.NotEmpty(t, w.Header().Get(httputil.HeaderRequestID))
}

func TestMiddleware_RequestID_Propagation(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.Header.Set(httputil.HeaderRequestID, "client-provided-id")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, "client-provided-id", w.Header().Get(httputil.HeaderRequestID))
}

func TestMiddleware_ApiVersion(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.NotEmpty(t, w.Header().Get("Api-Version"))
}

func TestHealth_Endpoint(t *testing.T) {
	repo := newMockTodoRepo()
	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
