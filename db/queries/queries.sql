-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id, username, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT id, username, password, created_at, updated_at
FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todos (title, description, completed, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id, title, description, completed, user_id, created_at, updated_at;

-- name: GetTodoByID :one
SELECT id, title, description, completed, user_id, created_at, updated_at
FROM todos
WHERE id = $1 AND user_id = $2;

-- name: ListTodos :many
SELECT id, title, description, completed, user_id, created_at, updated_at
FROM todos
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTodo :one
UPDATE todos
SET title = $1, description = $2, completed = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $4 AND user_id = $5
RETURNING id, title, description, completed, user_id, created_at, updated_at;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1 AND user_id = $2;

-- name: CountTodos :one
SELECT COUNT(*)
FROM todos
WHERE user_id = $1;
