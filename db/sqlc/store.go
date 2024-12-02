package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/achievagemilang/cheapbank/util"
)

type Store struct {
	*Queries
	db *sql.DB

}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccID int64    `json:"from_acc_id"`
	ToAccID   int64    `json:"to_acc_id"`
	Amount    float64    `json:"amount"`
	Currency  Currency `json:"currency"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccID,
			Amount: -arg.Amount,
			Currency: arg.Currency,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccID,
			Amount: arg.Amount,
			Currency: arg.Currency,
		})
		if err != nil {
			return err
		}

		account, err := q.GetAccountForUpdate(ctx, arg.FromAccID)
		if err != nil {
			return err
		}
		convertedAmount, err := util.ConvertCurrency(arg.Amount, string(arg.Currency), string(account.Currency))
		if err != nil {
			return err
		}
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccID,
			Amount: -convertedAmount,
		})
		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccID)
		if err != nil {
			return err
		}
		convertedAmount, err = util.ConvertCurrency(arg.Amount, string(arg.Currency), string(account2.Currency))
		if err != nil {
			return err
		}
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccID,
			Amount: convertedAmount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}