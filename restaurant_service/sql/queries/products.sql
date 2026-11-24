
-- name: GetProducts :many

SELECT
    id,
    name,
    price,
    created_at,
    updated_at,
    is_available
FROM products
ORDER BY id;

-- name: GetProduct :one

SELECT
    id,
    name,
    price,
    created_at,
    updated_at,
    is_available
FROM products
WHERE id = $1;

-- name: CreateProduct :one

INSERT INTO products (
    name,
    price
)
VALUES ($1, $2)
RETURNING *;