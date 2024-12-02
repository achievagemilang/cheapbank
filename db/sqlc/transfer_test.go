package db

import (
	"context"
	"testing"

	"github.com/achievagemilang/cheapbank/util"
	"github.com/stretchr/testify/require"
)


func createRandomTransfer(_ *testing.T, account1 Account, account2 Account) (Transfer, CreateTransferParams, error) {
	arg := CreateTransferParams{
		FromAccID: account1.ID,
		ToAccID:  account2.ID,
		Currency: Currency(util.RandomCurrency()),
		Amount:    util.RandomBalance(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	return transfer, arg, err
}

func TestCreateTransfer(t *testing.T) {
	account1, _, _ := createRandomAccount(t)
	account2, _, _ := createRandomAccount(t)
	transfer, arg, err := createRandomTransfer(t, account1, account2)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccID, transfer.FromAccID)
	require.Equal(t, arg.ToAccID, transfer.ToAccID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.Currency, transfer.Currency)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	account1, _, _ := createRandomAccount(t)
	account2, _, _ := createRandomAccount(t)
	transfer, _, err := createRandomTransfer(t, account1, account2)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer.Amount, transfer2.Amount)
	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.CreatedAt, transfer2.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	account1, _, _ := createRandomAccount(t)
	account2, _, _ := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
	}

	arg := ListTransfersParams{
		FromAccID: account1.ID,
		ToAccID: account2.ID,
		Limit:    5,
		Offset:   5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
