-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    amount,
    status,
    method,
    transaction_id,
    error_message
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, order_id, amount, status, method, transaction_id, error_message, created_at;
