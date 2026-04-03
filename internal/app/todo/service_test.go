package todo_test

import (
	"context"
	"errors"
	"testing"

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

// Tests

func TestService_Create(t *testing.T) {
	t.Run("creates todo successfully", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		got, err := svc.Create(context.Background(), apptodo.CreateCommand{
			Title:       "Buy groceries",
			Description: "milk, eggs",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, "Buy groceries", got.Title)
		assert.False(t, got.Completed)
	})

	t.Run("empty title returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		_, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: ""})
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrTitleRequired))
	})
}

func TestService_FindByID(t *testing.T) {
	t.Run("finds existing todo", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Test"})
		require.NoError(t, err)

		got, err := svc.FindByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
	})

	t.Run("not found returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		_, err := svc.FindByID(context.Background(), "nonexistent")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})
}

func TestService_FindAll(t *testing.T) {
	t.Run("returns empty list when no todos", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		result, err := svc.FindAll(context.Background(), apptodo.ListQuery{})
		require.NoError(t, err)
		assert.Empty(t, result.Items)
		assert.Equal(t, int64(0), result.TotalCount)
	})

	t.Run("returns all todos", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		_, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Todo 1"})
		require.NoError(t, err)
		_, err = svc.Create(context.Background(), apptodo.CreateCommand{Title: "Todo 2"})
		require.NoError(t, err)

		result, err := svc.FindAll(context.Background(), apptodo.ListQuery{})
		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.TotalCount)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("updates title", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Old"})
		require.NoError(t, err)

		newTitle := "New Title"
		got, err := svc.Update(context.Background(), created.ID, apptodo.UpdateCommand{
			Title: &newTitle,
		})
		require.NoError(t, err)
		assert.Equal(t, "New Title", got.Title)
	})

	t.Run("updates description", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Title", Description: "Old desc"})
		require.NoError(t, err)

		newDesc := "New desc"
		got, err := svc.Update(context.Background(), created.ID, apptodo.UpdateCommand{
			Description: &newDesc,
		})
		require.NoError(t, err)
		assert.Equal(t, "New desc", got.Description)
	})

	t.Run("not found returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		newTitle := "Title"
		_, err := svc.Update(context.Background(), "nonexistent", apptodo.UpdateCommand{Title: &newTitle})
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})

	t.Run("empty title update returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Title"})
		require.NoError(t, err)

		emptyTitle := ""
		_, err = svc.Update(context.Background(), created.ID, apptodo.UpdateCommand{Title: &emptyTitle})
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrTitleRequired))
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("deletes existing todo", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "To delete"})
		require.NoError(t, err)

		require.NoError(t, svc.Delete(context.Background(), created.ID))

		_, err = svc.FindByID(context.Background(), created.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("deleting non-existent returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		err := svc.Delete(context.Background(), "nonexistent")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})
}

func TestService_Complete(t *testing.T) {
	t.Run("completes a todo", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Test"})
		require.NoError(t, err)

		got, err := svc.Complete(context.Background(), created.ID)
		require.NoError(t, err)
		assert.True(t, got.Completed)
	})

	t.Run("completing already done returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		created, err := svc.Create(context.Background(), apptodo.CreateCommand{Title: "Test"})
		require.NoError(t, err)

		_, err = svc.Complete(context.Background(), created.ID)
		require.NoError(t, err)

		_, err = svc.Complete(context.Background(), created.ID)
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrAlreadyDone))
	})

	t.Run("completing non-existent returns error", func(t *testing.T) {
		repo := newMockRepo()
		svc := apptodo.NewService(repo)

		_, err := svc.Complete(context.Background(), "nonexistent")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})
}
