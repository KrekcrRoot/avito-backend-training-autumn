
-- name: GetUsers :many
SELECT
    id,
    name,
    email,
    password,
    created_at,
    updated_at
FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE email = $1;