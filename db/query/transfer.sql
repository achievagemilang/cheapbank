-- name: CreateTransfer :one
INSERT INTO transfers (
    from_acc_id,
    to_acc_id,
    amount,
    currency
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers 
WHERE id = $1
LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE
    from_acc_id = $1 OR
    to_acc_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;
