package db

import (
	"context"
	"testing"

	"github.com/achievagemilang/cheapbank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(_ *testing.T, account Account) (Entry, CreateEntryParams, error) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
		Currency:  account.Currency,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	return entry, arg, err
}


func TestCreateEntry(t *testing.T) {
	account, _, _ := createRandomAccount(t)
	entry, arg, err := createRandomEntry(t, account)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	account, _, _ := createRandomAccount(t)
	entry, _, err := createRandomEntry(t, account)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry.Amount, entry2.Amount)
	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.CreatedAt, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account, _, _ := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t,account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}