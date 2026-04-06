package httphandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/eventbus"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/httphandler"
	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
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
	bus := eventbus.NewInMemoryEventBus()
	svc := apptodo.NewService(repo, bus)
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

func TestHandler_Get(t *testing.T) {
	repo := newMockTodoRepo()
	id := uuid.New().String()
	td, _ := domain.New(id, "Test Todo", "some desc")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodGet, "/todos/"+id, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, id, resp["id"])
}

func TestHandler_Update(t *testing.T) {
	repo := newMockTodoRepo()
	id := uuid.New().String()
	td, _ := domain.New(id, "Old Title", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	body := `{"title":"New Title"}`
	r := httptest.NewRequest(http.MethodPatch, "/todos/"+id, bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "New Title", resp["title"])
}

func TestHandler_Delete(t *testing.T) {
	repo := newMockTodoRepo()
	id := uuid.New().String()
	td, _ := domain.New(id, "To delete", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodDelete, "/todos/"+id, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestHandler_Complete(t *testing.T) {
	repo := newMockTodoRepo()
	id := uuid.New().String()
	td, _ := domain.New(id, "Test", "")
	_ = repo.Save(context.Background(), td)

	router := setupRouter(repo)

	r := httptest.NewRequest(http.MethodPost, "/todos/"+id+":complete", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, true, resp["completed"])
}
