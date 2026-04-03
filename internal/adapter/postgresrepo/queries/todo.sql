-- name: CreateTodo :one
INSERT INTO todos (id, title, description, completed, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1;

-- name: ListTodos :many
SELECT * FROM todos
WHERE ($1::varchar = '' OR id > $1)
ORDER BY id ASC
LIMIT $2;

-- name: UpdateTodo :one
UPDATE todos
SET title = $2, description = $3, completed = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;

-- name: CountTodos :one
SELECT COUNT(*) FROM todos;
