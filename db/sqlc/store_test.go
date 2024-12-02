package db

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/achievagemilang/cheapbank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1, _, _ := createRandomAccount(t)
	account2, _, _ := createRandomAccount(t)

	n := 5
	amount := float64(10)
	currency := Currency(util.RandomCurrency())

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	var wg sync.WaitGroup
	wg.Add(n)

	type ctxKey string

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			defer wg.Done()
			txKey := ctxKey(fmt.Sprintf("tx %d", i+1))
			var ctx = context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccID: account1.ID,
				ToAccID:   account2.ID,
				Amount:    amount,
				Currency:  currency,
			})

			errs <- err
			results <- result
		}()
	}

	wg.Wait()
	close(errs)
	close(results)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccID)
		require.Equal(t, account2.ID, transfer.ToAccID)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, currency, transfer.Currency)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.Equal(t, currency, fromEntry.Currency)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.Equal(t, currency, toEntry.Currency)
		require.NotZero(t, toEntry.ID)

		// TODO: Additional checks like account balances

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		diff1, _ = util.ConvertCurrency(diff1, string(account1.Currency), string(currency))
		diff2, _ = util.ConvertCurrency(diff2, string(account2.Currency), string(currency))
		margin := 0.00000000001
		require.InEpsilon(t, diff1, diff2, margin)
		require.True(t, diff1 > 0)
	}
}
