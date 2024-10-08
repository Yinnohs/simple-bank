// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account.sql

package db

import (
	"context"
)

const addToAccountBalance = `-- name: AddToAccountBalance :one
UPDATE accounts 
SET balance = balance + $1
WHERE id = $2
RETURNING id, owner, currency, balance, created_at
`

type AddToAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddToAccountBalance(ctx context.Context, arg AddToAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addToAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Currency,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  owner, currency, balance
) VALUES (
  $1, $2, $3
) RETURNING id, owner, currency, balance, created_at
`

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Currency, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Currency,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccountById = `-- name: DeleteAccountById :exec
DELETE FROM accounts 
WHERE id = $1
`

func (q *Queries) DeleteAccountById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccountById, id)
	return err
}

const findAllAccounts = `-- name: FindAllAccounts :many

SELECT id, owner, currency, balance, created_at 
FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2
`

type FindAllAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// FOR UPDATE meake a select for transactional queries, it locks the row until the transaction is finished
func (q *Queries) FindAllAccounts(ctx context.Context, arg FindAllAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, findAllAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Currency,
			&i.Balance,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAccountById = `-- name: GetAccountById :one
SELECT id, owner, currency, balance, created_at 
FROM accounts 
WHERE id = $1 
LIMIT 1
`

func (q *Queries) GetAccountById(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Currency,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountForUpdateById = `-- name: GetAccountForUpdateById :one
SELECT id, owner, currency, balance, created_at 
FROM accounts 
WHERE id = $1 
LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdateById(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdateById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Currency,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const updateAccountBalance = `-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING id, owner, currency, balance, created_at
`

type UpdateAccountBalanceParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccountBalance, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Currency,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}
