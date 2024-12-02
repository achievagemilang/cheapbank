package db

import (
	"context"
	"testing"

	"github.com/achievagemilang/cheapbank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(_ *testing.T) (Account, CreateAccountParams, error) {
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Currency: Currency(util.RandomCurrency()),
		Balance: util.RandomBalance(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	return account, arg, err
}

func TestCreateAccount(t *testing.T){
	account, arg, err:= createRandomAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T){
	account, _, err := createRandomAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.Equal(t, account.Balance, account2.Balance)
	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)
}

func TestDeleteAccount(t *testing.T){
	account, _, err := createRandomAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	err = testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T){
	var lastAccount Account
	for i := 0; i < 10; i++ {
		account, _, err := createRandomAccount(t)
		require.NoError(t, err)
		require.NotEmpty(t, account)

		require.NotZero(t, account.ID)
		require.NotZero(t, account.CreatedAt)

		if lastAccount.ID != 0 {
			require.NotEqual(t, lastAccount.ID, account.ID)
		}

		lastAccount = account
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T){
	account, _, err := createRandomAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	arg := UpdateAccountParams{}

	arg.Balance = util.RandomBalance()
	arg.ID = account.ID

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)
}