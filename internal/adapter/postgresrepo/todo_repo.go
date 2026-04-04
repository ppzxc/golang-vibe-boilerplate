// Package postgresrepo implements the domain Repository ports using PostgreSQL.
package postgresrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

// TodoRepository implements domain/todo.Repository using PostgreSQL.
type TodoRepository struct {
	db *sql.DB
}

// NewTodoRepository creates a new TodoRepository.
func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Save inserts a new Todo into the database.
func (r *TodoRepository) Save(ctx context.Context, todo *domain.Todo) error {
	q := `INSERT INTO todos (id, title, description, completed, created_at, updated_at)
          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, q,
		todo.ID, todo.Title, todo.Description, todo.Completed,
		todo.CreatedAt, todo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("postgresrepo.TodoRepository.Save: %w", err)
	}
	return nil
}

// FindByID retrieves a Todo by its ID.
// Returns domain.ErrNotFound if the Todo does not exist.
func (r *TodoRepository) FindByID(ctx context.Context, id string) (*domain.Todo, error) {
	q := `SELECT id, title, description, completed, created_at, updated_at
          FROM todos WHERE id = $1`
	row := r.db.QueryRowContext(ctx, q, id)
	todo, err := scanTodo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("postgresrepo.TodoRepository.FindByID: %w", err)
	}
	return todo, nil
}

// FindAll retrieves a page of Todos using cursor-based pagination.
func (r *TodoRepository) FindAll(ctx context.Context, cursor string, pageSize int) ([]*domain.Todo, string, error) {
	q := `SELECT id, title, description, completed, created_at, updated_at
          FROM todos
          WHERE ($1::varchar = '' OR id > $1)
          ORDER BY id ASC
          LIMIT $2`
	rows, err := r.db.QueryContext(ctx, q, cursor, pageSize+1)
	if err != nil {
		return nil, "", fmt.Errorf("postgresrepo.TodoRepository.FindAll: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var todos []*domain.Todo
	for rows.Next() {
		todo, err := scanTodoRow(rows)
		if err != nil {
			return nil, "", fmt.Errorf("postgresrepo.TodoRepository.FindAll: %w", err)
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("postgresrepo.TodoRepository.FindAll: %w", err)
	}

	nextCursor := ""
	if len(todos) > pageSize {
		nextCursor = todos[pageSize-1].ID
		todos = todos[:pageSize]
	}

	return todos, nextCursor, nil
}

// Update retrieves a Todo, applies fn to it, and saves the result.
// This is the UpdateFn pattern for transactional updates.
func (r *TodoRepository) Update(ctx context.Context, id string, fn func(*domain.Todo) (*domain.Todo, error)) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgresrepo.TodoRepository.Update: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	q := `SELECT id, title, description, completed, created_at, updated_at
          FROM todos WHERE id = $1 FOR UPDATE`
	row := tx.QueryRowContext(ctx, q, id)
	todo, err := scanTodo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("postgresrepo.TodoRepository.Update: %w", err)
	}

	updated, err := fn(todo)
	if err != nil {
		return err
	}

	uq := `UPDATE todos
           SET title = $2, description = $3, completed = $4, updated_at = $5
           WHERE id = $1`
	if _, err := tx.ExecContext(ctx, uq,
		updated.ID, updated.Title, updated.Description,
		updated.Completed, updated.UpdatedAt,
	); err != nil {
		return fmt.Errorf("postgresrepo.TodoRepository.Update: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("postgresrepo.TodoRepository.Update: commit: %w", err)
	}
	return nil
}

// Delete removes a Todo by ID.
func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM todos WHERE id = $1`
	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("postgresrepo.TodoRepository.Delete: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgresrepo.TodoRepository.Delete: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Count returns the total number of Todos.
func (r *TodoRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM todos`).Scan(&count); err != nil {
		return 0, fmt.Errorf("postgresrepo.TodoRepository.Count: %w", err)
	}
	return count, nil
}

// Open opens a PostgreSQL connection with the given DSN.
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgresrepo.Open: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTodo(s scanner) (*domain.Todo, error) {
	var t domain.Todo
	err := s.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func scanTodoRow(rows *sql.Rows) (*domain.Todo, error) {
	var t domain.Todo
	err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
