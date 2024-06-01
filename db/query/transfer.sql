-- name: CreateTransfers :one
INSERT INTO transfers (
  from_account_id, 
  to_account_id,
  amount, 
  status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;
