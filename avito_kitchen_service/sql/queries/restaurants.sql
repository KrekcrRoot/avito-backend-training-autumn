
-- name: GetRestaurants :many
SELECT
    id,
    name,
    base_url,
    api_key,
    created_at,
    updated_at
FROM restaurants
ORDER BY id;

-- name: GetRestaurantByID :one
SELECT
    id,
    name,
    base_url,
    api_key,
    created_at,
    updated_at
FROM restaurants
WHERE id = $1;

-- name: GetRestaurantByAPIKey :one
SELECT id,
       name,
       base_url,
       api_key,
       created_at,
       updated_at
FROM restaurants
WHERE api_key = $1;