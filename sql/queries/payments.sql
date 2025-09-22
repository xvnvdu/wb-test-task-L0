-- name: CreatePayment :exec
INSERT INTO payments (
    order_uid, 
    transaction,
    request_id,
    currency,
    provider,
    amount,
    payment_dt,
    bank,
    delivery_cost,
    goods_total,
    custom_fee
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetPayment :many
SELECT * FROM payments;

-- name: GetSpecificPayment :one
SELECT * FROM payments WHERE order_uid = $1;
