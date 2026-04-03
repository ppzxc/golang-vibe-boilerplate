//go:build integration

package postgresrepo_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/postgresrepo"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPostgres(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port())

	db, err := postgresrepo.Open(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Run migration
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS todos (
			id          VARCHAR(36)  PRIMARY KEY,
			title       VARCHAR(255) NOT NULL,
			description TEXT         NOT NULL DEFAULT '',
			completed   BOOLEAN      NOT NULL DEFAULT FALSE,
			created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	return db
}

func TestTodoRepository_Save_Integration(t *testing.T) {
	db := setupPostgres(t)
	repo := postgresrepo.NewTodoRepository(db)
	ctx := context.Background()

	td, err := domain.New("id-1", "Buy milk", "2% milk")
	require.NoError(t, err)

	err = repo.Save(ctx, td)
	require.NoError(t, err)

	got, err := repo.FindByID(ctx, "id-1")
	require.NoError(t, err)
	assert.Equal(t, "id-1", got.ID)
	assert.Equal(t, "Buy milk", got.Title)
	assert.False(t, got.Completed)
}

func TestTodoRepository_FindByID_NotFound_Integration(t *testing.T) {
	db := setupPostgres(t)
	repo := postgresrepo.NewTodoRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, "nonexistent")
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTodoRepository_Update_Integration(t *testing.T) {
	db := setupPostgres(t)
	repo := postgresrepo.NewTodoRepository(db)
	ctx := context.Background()

	td, _ := domain.New("id-2", "Original", "")
	require.NoError(t, repo.Save(ctx, td))

	err := repo.Update(ctx, "id-2", func(t *domain.Todo) (*domain.Todo, error) {
		return t, t.UpdateTitle("Updated")
	})
	require.NoError(t, err)

	got, _ := repo.FindByID(ctx, "id-2")
	assert.Equal(t, "Updated", got.Title)
}

func TestTodoRepository_Delete_Integration(t *testing.T) {
	db := setupPostgres(t)
	repo := postgresrepo.NewTodoRepository(db)
	ctx := context.Background()

	td, _ := domain.New("id-3", "Delete me", "")
	require.NoError(t, repo.Save(ctx, td))

	require.NoError(t, repo.Delete(ctx, "id-3"))

	_, err := repo.FindByID(ctx, "id-3")
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTodoRepository_FindAll_Integration(t *testing.T) {
	db := setupPostgres(t)
	repo := postgresrepo.NewTodoRepository(db)
	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		td, _ := domain.New(fmt.Sprintf("id-%d", i), fmt.Sprintf("Todo %d", i), "")
		require.NoError(t, repo.Save(ctx, td))
	}

	todos, nextCursor, err := repo.FindAll(ctx, "", 10)
	require.NoError(t, err)
	assert.Len(t, todos, 3)
	assert.Empty(t, nextCursor)
}

func TestTodoRepository_Count_Integration(t *testing.T) {
	db := setupPostgres(t)
	repo := postgresrepo.NewTodoRepository(db)
	ctx := context.Background()

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	td, _ := domain.New("id-1", "Test", "")
	require.NoError(t, repo.Save(ctx, td))

	count, err = repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
