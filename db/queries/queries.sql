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

-- URL Shortener Queries
-- name: CreateShortURL :one
INSERT INTO short_urls (code, original_url, alias, clicks, is_public, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, code, original_url, alias, clicks, is_public, expires_at, created_at, updated_at;

-- name: GetShortURLByCode :one
SELECT id, code, original_url, alias, clicks, is_public, expires_at, created_at, updated_at
FROM short_urls
WHERE code = $1;

-- name: IncrementShortURLClicks :exec
UPDATE short_urls
SET clicks = clicks + 1, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateURLClick :exec
INSERT INTO url_clicks (short_url_id, referrer, user_agent, ip_address)
VALUES ($1, $2, $3, $4);

-- name: DeleteExpiredShortURLs :exec
DELETE FROM short_urls
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;

-- Pastebin Queries
-- name: CreatePaste :one
INSERT INTO pastes (id, title, content, syntax, is_public, is_compressed, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, content, syntax, is_public, is_compressed, expires_at, created_at, updated_at;

-- name: GetPasteByID :one
SELECT id, title, content, syntax, is_public, is_compressed, expires_at, created_at, updated_at
FROM pastes
WHERE id = $1;

-- name: DeletePaste :exec
DELETE FROM pastes
WHERE id = $1;

-- name: ListRecentPastes :many
SELECT id, title, content, syntax, is_public, is_compressed, expires_at, created_at, updated_at
FROM pastes
WHERE is_public = true
ORDER BY created_at DESC
LIMIT $1;

-- name: DeleteExpiredPastes :exec
DELETE FROM pastes
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;

-- QR Code Queries
-- name: CreateQRCode :one
INSERT INTO qr_codes (id, text, format, size, image_data)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, text, format, size, created_at;

-- name: GetQRCodeByID :one
SELECT id, text, format, size, image_data, created_at
FROM qr_codes
WHERE id = $1;
