
-- name: GetOrders :many

SELECT
    id,
    kitchen_order_id,
    status,
    created_at,
    updated_at
FROM orders
ORDER BY id;

-- name: GetOrder :one
SELECT
    id,
    kitchen_order_id,
    status,
    created_at,
    updated_at
FROM orders
WHERE id = $1;

-- name: GetOrderItems :many
SELECT
    order_id,
    product_id,
    quantity,
    price
FROM order_items
WHERE order_id = $1;

-- name: CreateOrder :one

INSERT INTO orders (
    kitchen_order_id
)
VALUES ($1)
RETURNING *;

-- name: AddOrderItem :one

INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;