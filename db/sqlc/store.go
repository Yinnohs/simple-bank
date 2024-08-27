package db

import (
	"context"
	"database/sql"
	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// function that executes a trasanccions wrapper

func (store *Store) executeTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal(err)
		return err
	}

	query := New(tx)

	err = fn(query)
	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			log.Fatal("transaction error: " + err.Error() + " rollback error: " + rbError.Error())
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.executeTransaction(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.GetAccountById(ctx, args.FromAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.GetAccountById(ctx, args.ToAccountID)
		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				args.FromAccountID,
				-args.Amount,
				args.ToAccountID,
				args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx,
				q,
				args.ToAccountID,
				args.Amount,
				args.FromAccountID,
				-args.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountIDFrom int64,
	amountToSubstract int64,
	accountIDTo int64,
	amountToAdd int64,
) (accountFrom Account, accountTo Account, err error) {

	accountFrom, err = q.AddToAccountBalance(ctx, AddToAccountBalanceParams{
		ID:     accountIDFrom,
		Amount: amountToSubstract,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	accountTo, err = q.AddToAccountBalance(ctx, AddToAccountBalanceParams{
		ID:     accountIDTo,
		Amount: amountToAdd,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	return
}
