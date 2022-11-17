CREATE TABLE IF NOT EXISTS carts (
    user_id uuid,
    good_id uuid,
    count int,
    created_at timestamptz not null default current_timestamp
);

---- create above / drop below ----

DROP TABLE IF EXISTS cartsж
