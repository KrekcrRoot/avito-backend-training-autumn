
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    kitchen_order_id BIGINT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'new' CHECK (
        status IN (
               'new',
               'accepted',
               'cooking',
               'ready',
               'rejected',
               'canceled'
            )
        ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);