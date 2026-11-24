
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    restaurant_order_id BIGINT,
    status TEXT NOT NULL DEFAULT 'new',
    total_price NUMERIC(10, 2) NOT NULL CHECK (total_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);