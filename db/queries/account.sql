-- name: CreateAccount :one
INSERT INTO accounts (
  owner, currency, balance
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccountById :one
SELECT * 
FROM accounts 
WHERE id = $1 
LIMIT 1;

-- name: GetAccountForUpdateById :one
SELECT * 
FROM accounts 
WHERE id = $1 
LIMIT 1
FOR NO KEY UPDATE;

-- FOR UPDATE meake a select for transactional queries, it locks the row until the transaction is finished

-- name: FindAllAccounts :many
SELECT * 
FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddToAccountBalance :one
UPDATE accounts 
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccountById :exec
DELETE FROM accounts 
WHERE id = $1;
