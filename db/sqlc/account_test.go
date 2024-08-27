package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yinnohs/simple-bank/util"
)

func createRandomAccount(params CreateAccountParams) (Account, error) {

	account, err := testQueries.CreateAccount(context.Background(), params)
	if err != nil {
		return Account{}, err
	}

	return account, nil

}

func TestCreateAccount(t *testing.T) {
	// when
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := createRandomAccount(params)
	//then
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	//given
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	createdAccount, _ := createRandomAccount(params)

	//when
	account, err := testQueries.GetAccountById(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)
	require.Equal(t, createdAccount.ID, account.ID)
	require.WithinDuration(t, createdAccount.CreatedAt, account.CreatedAt, time.Second)

}

func TestUpdateAccountBalance(t *testing.T) {
	createParams := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	createdAccount, _ := createRandomAccount(createParams)

	updateParams := UpdateAccountBalanceParams{
		ID:      createdAccount.ID,
		Balance: 2000,
	}

	updatedAccount, err := testQueries.UpdateAccountBalance(context.Background(), updateParams)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.WithinDuration(t, createdAccount.CreatedAt, updatedAccount.CreatedAt, time.Second)
	require.Equal(t, updateParams.Balance, updatedAccount.Balance)

}
