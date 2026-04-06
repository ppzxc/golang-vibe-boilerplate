package todo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/eventbus"
	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepo is a simple in-memory mock of domain.Repository.
type mockRepo struct {
	todos   map[string]*domain.Todo
	saveErr error
	findErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{todos: make(map[string]*domain.Todo)}
}

func (m *mockRepo) Save(_ context.Context, t *domain.Todo) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.todos[t.ID] = t
	return nil
}

func (m *mockRepo) FindByID(_ context.Context, id string) (*domain.Todo, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	t, ok := m.todos[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (m *mockRepo) FindAll(_ context.Context, _ string, _ int) ([]*domain.Todo, string, error) {
	var result []*domain.Todo
	for _, t := range m.todos {
		cp := *t
		result = append(result, &cp)
	}
	return result, "", nil
}

func (m *mockRepo) Update(_ context.Context, id string, fn func(*domain.Todo) (*domain.Todo, error)) error {
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

func (m *mockRepo) Delete(_ context.Context, id string) error {
	if _, ok := m.todos[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.todos, id)
	return nil
}

func (m *mockRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.todos)), nil
}

func setupService(repo domain.Repository) *apptodo.Service {
	bus := eventbus.NewInMemoryEventBus()
	return apptodo.NewService(repo, bus)
}

// Tests

func TestService_Create(t *testing.T) {
	t.Run("creates todo successfully", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		got, err := svc.Create(context.Background(), "Buy groceries", "milk, eggs")

		require.NoError(t, err)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, "Buy groceries", got.Title)
		assert.False(t, got.Completed)
	})

	t.Run("empty title returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		_, err := svc.Create(context.Background(), "", "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrTitleRequired))
	})
}

func TestService_FindByID(t *testing.T) {
	t.Run("finds existing todo", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		created, err := svc.Create(context.Background(), "Test", "")
		require.NoError(t, err)

		got, err := svc.FindByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
	})

	t.Run("not found returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		_, err := svc.FindByID(context.Background(), "nonexistent")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})
}

func TestService_FindAll(t *testing.T) {
	t.Run("returns empty list when no todos", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		items, _, count, err := svc.FindAll(context.Background(), "", 0)
		require.NoError(t, err)
		assert.Empty(t, items)
		assert.Equal(t, int64(0), count)
	})

	t.Run("returns all todos", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		_, err := svc.Create(context.Background(), "Todo 1", "")
		require.NoError(t, err)
		_, err = svc.Create(context.Background(), "Todo 2", "")
		require.NoError(t, err)

		items, _, count, err := svc.FindAll(context.Background(), "", 0)
		require.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, int64(2), count)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("updates title", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		created, err := svc.Create(context.Background(), "Old", "")
		require.NoError(t, err)

		newTitle := "New Title"
		got, err := svc.Update(context.Background(), created.ID, &newTitle, nil)
		require.NoError(t, err)
		assert.Equal(t, "New Title", got.Title)
	})

	t.Run("updates description", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		created, err := svc.Create(context.Background(), "Title", "Old desc")
		require.NoError(t, err)

		newDesc := "New desc"
		got, err := svc.Update(context.Background(), created.ID, nil, &newDesc)
		require.NoError(t, err)
		assert.Equal(t, "New desc", got.Description)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("deletes existing todo", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		created, err := svc.Create(context.Background(), "To delete", "")
		require.NoError(t, err)

		require.NoError(t, svc.Delete(context.Background(), created.ID))

		_, err = svc.FindByID(context.Background(), created.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestService_Complete(t *testing.T) {
	t.Run("completes a todo", func(t *testing.T) {
		repo := newMockRepo()
		svc := setupService(repo)

		created, err := svc.Create(context.Background(), "Test", "")
		require.NoError(t, err)

		got, err := svc.Complete(context.Background(), created.ID)
		require.NoError(t, err)
		assert.True(t, got.Completed)
	})
}
