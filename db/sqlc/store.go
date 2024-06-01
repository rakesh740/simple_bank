package db

import (
	"context"
	"database/sql"
	"fmt"
)

type IStore interface {
	Querier
	Transfertx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error)
}

// store provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) IStore {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollErr := tx.Rollback(); rollErr != nil {
			return fmt.Errorf("tx Err: %v rollback err: %v", err, rollErr)
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

type TransferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) Transfertx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {

	var result TransferTxResults
	err := store.execTx(ctx, func(q *Queries) (err error) {
		if result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}); err != nil {
			return err
		}

		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}); err != nil {
			return err
		}

		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}); err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {

			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
			if err != nil {
				return err
			}

		} else {

			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context,
	q *Queries,
	accountID1, accountID2,
	amount1, amount2 int64) (account1, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
