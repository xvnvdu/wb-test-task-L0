-- name: CreateOrder :exec
INSERT INTO orders (
    order_uid, 
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetOrders :many
SELECT * FROM orders;

-- name: GetSpecificOrder :one
SELECT * FROM orders WHERE order_uid = $1;
