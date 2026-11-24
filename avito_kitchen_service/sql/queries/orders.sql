
-- name: GetOrderByID :one
SELECT
    id, restaurant_id, user_id, restaurant_order_id, status, total_price, created_at
FROM orders
WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT
    id, restaurant_id, user_id, restaurant_order_id, status, total_price, created_at
FROM orders
WHERE user_id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    restaurant_id, user_id, total_price
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddOrderSnapshotItem :one
INSERT INTO order_items (
    order_id, restaurant_product_id, name, price, quantity
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: SetRestaurantOrderID :one
UPDATE orders
SET restaurant_order_id = $2
WHERE id = $1
RETURNING *;

-- name: GetOrderItemsByOrderID :many
SELECT
    order_id,
    restaurant_product_id,
    name,
    price,
    quantity
FROM order_items
WHERE order_id = $1
ORDER BY restaurant_product_id;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $3
WHERE restaurant_id = $1
  AND restaurant_order_id = $2
RETURNING *;