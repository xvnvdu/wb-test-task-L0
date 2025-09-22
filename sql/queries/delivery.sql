-- name: CreateDelivery :exec
INSERT INTO delivery (
    order_uid, 
    name,
    phone,
    zip,
    city,
    address,
    region,
    email
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDelivery :many
SELECT * FROM delivery;

-- name: GetSpecificDelivery :one
SELECT * FROM delivery WHERE order_uid = $1;
