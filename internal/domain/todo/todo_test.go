package todo_test

import (
	"testing"

	"github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     error
	}{
		{name: "valid", title: "Buy groceries", description: "milk, eggs"},
		{name: "empty title", title: "", wantErr: todo.ErrTitleRequired},
		{name: "no description ok", title: "Title only", description: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := todo.New("id-1", tc.title, tc.description)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "id-1", got.ID)
			assert.Equal(t, tc.title, got.Title)
			assert.False(t, got.Completed)
			assert.NotZero(t, got.CreatedAt)
		})
	}
}

func TestTodo_Complete(t *testing.T) {
	t.Run("completes successfully", func(t *testing.T) {
		td, err := todo.New("id-1", "Test", "")
		require.NoError(t, err)
		require.NoError(t, td.Complete())
		assert.True(t, td.Completed)
	})

	t.Run("already completed returns error", func(t *testing.T) {
		td, err := todo.New("id-1", "Test", "")
		require.NoError(t, err)
		require.NoError(t, td.Complete())
		assert.ErrorIs(t, td.Complete(), todo.ErrAlreadyDone)
	})
}

func TestTodo_UpdateTitle(t *testing.T) {
	t.Run("updates title", func(t *testing.T) {
		td, _ := todo.New("id-1", "Old", "")
		require.NoError(t, td.UpdateTitle("New"))
		assert.Equal(t, "New", td.Title)
	})

	t.Run("empty title returns error", func(t *testing.T) {
		td, _ := todo.New("id-1", "Title", "")
		assert.ErrorIs(t, td.UpdateTitle(""), todo.ErrTitleRequired)
	})
}

func TestTodo_UpdateDescription(t *testing.T) {
	t.Run("updates description", func(t *testing.T) {
		td, _ := todo.New("id-1", "Title", "old desc")
		td.UpdateDescription("new desc")
		assert.Equal(t, "new desc", td.Description)
	})

	t.Run("clears description", func(t *testing.T) {
		td, _ := todo.New("id-1", "Title", "old desc")
		td.UpdateDescription("")
		assert.Equal(t, "", td.Description)
	})
}

func TestTodo_String(t *testing.T) {
	td, _ := todo.New("id-1", "Test", "")
	s := td.String()
	assert.Contains(t, s, "id-1")
	assert.Contains(t, s, "Test")
}
