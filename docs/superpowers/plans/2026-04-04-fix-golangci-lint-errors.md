# golangci-lint Error Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix all golangci-lint errors discovered in CI and improve code quality.

**Architecture:** Rename domain events, fix non-idiomatic parameter names, delete unused code, and wrap external errors.

**Tech Stack:** Go, golangci-lint, oapi-codegen, gofumpt

---

### Task 1: Rename Domain Events

**Files:**
- Modify: `internal/domain/todo/todo.go`
- Modify: `internal/adapter/eventbus/inmem.go`

- [ ] **Step 1: Rename types and methods in `internal/domain/todo/todo.go`**

```go
// In internal/domain/todo/todo.go

// Completed is a domain event.
type Completed struct {
	ID   string
	When time.Time
}

func (e Completed) OccurredAt() time.Time { return e.When }

// Created is a domain event.
type Created struct {
	ID   string
	When time.Time
}

func (e Created) OccurredAt() time.Time { return e.When }

// Updated is a domain event.
type Updated struct {
	ID   string
	When time.Time
}

func (e Updated) OccurredAt() time.Time { return e.When }
```

- [ ] **Step 2: Update internal references in `internal/domain/todo/todo.go`**

Update `TodoCreated{}`, `TodoCompleted{}`, `TodoUpdated{}` to `Created{}`, `Completed{}`, `Updated{}`.

- [ ] **Step 3: Update references in `internal/adapter/eventbus/inmem.go`**

Update `domain.TodoCreated` to `domain.Created`, etc., and update the returned strings.

---

### Task 2: Rename 'todoId' to 'todoID'

**Files:**
- Modify: `api/openapi.yaml`
- Modify: `internal/adapter/httphandler/todo_handler.go`
- Regenerate: `internal/adapter/httphandler/api.gen.go`

- [ ] **Step 1: Update `api/openapi.yaml`**

Change all occurrences of `todoId` to `todoID` in path parameters.

- [ ] **Step 2: Run code generation**

Run: `make generate`

- [ ] **Step 3: Update `internal/adapter/httphandler/todo_handler.go`**

Update method signatures and local variables from `todoId` to `todoID`.

---

### Task 3: Delete Unused Code in httphandler

**Files:**
- Modify: `internal/adapter/httphandler/error.go`
- Modify: `internal/adapter/httphandler/request.go`
- Modify: `internal/adapter/httphandler/response.go`

- [ ] **Step 1: Delete `writeError` from `error.go`**
- [ ] **Step 2: Delete unused types from `request.go`**
- [ ] **Step 3: Delete unused types and functions from `response.go`**

---

### Task 4: Wrap Errors in postgresrepo

**Files:**
- Modify: `internal/adapter/postgresrepo/todo_repo.go`

- [ ] **Step 1: Wrap errors in `scanTodo` and `scanTodoRow`**

```go
func scanTodo(s scanner) (*domain.Todo, error) {
	var t domain.Todo
	err := s.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan todo: %w", err)
	}
	return &t, nil
}

func scanTodoRow(rows *sql.Rows) (*domain.Todo, error) {
	var t domain.Todo
	err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan todo row: %w", err)
	}
	return &t, nil
}
```

---

### Task 5: Formatting and Final Verification

- [ ] **Step 1: Run `gofumpt -w .`**
- [ ] **Step 2: Verify build and run lint**

Run: `go build ./...`
Run: `make lint` (or `golangci-lint run ./...` if available)
