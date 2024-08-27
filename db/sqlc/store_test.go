package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yinnohs/simple-bank/util"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testDB)
	params1 := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	params2 := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account1, _ := createRandomAccount(params1)
	account2, _ := createRandomAccount(params2)

	// run it in concurrent go routines
	loops := 5

	// when
	errs := make(chan error)
	results := make(chan TransferTxResult)
	var amount int64 = 10

	for i := 0; i < loops; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	//then
	for i := 0; i < loops; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		// checks entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: should check the balance next

		fromAccount := result.FromAccount
		require.NotNil(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotNil(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diffFromAccount := account1.Balance - fromAccount.Balance
		diffToAcount := toAccount.Balance - account2.Balance

		require.Equal(t, diffFromAccount, diffToAcount)
		require.True(t, diffFromAccount > 0)
		require.True(t, diffFromAccount%amount == 0)
		require.True(t, diffToAcount > 0)
		require.True(t, diffToAcount%amount == 0)

		resultDiffFromAmount := int(diffFromAccount / amount)
		require.True(t, resultDiffFromAmount >= 1 && resultDiffFromAmount <= loops)
		require.NotContains(t, existed, resultDiffFromAmount)
		existed[resultDiffFromAmount] = true
	}

	// check the final updated balance of th transaction

	updatedFromAccount, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	// account 1 is fromAccount and account 2 is toAccount
	require.Equal(t, account1.Balance-int64(loops)*amount, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(loops)*amount, updatedToAccount.Balance)
}

func TestTranferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	params1 := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	params2 := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account1, _ := createRandomAccount(params1)
	account2, _ := createRandomAccount(params2)

	// run it in concurrent go routines
	loops := 10

	// when
	errs := make(chan error)
	var amount int64 = 10

	for i := 0; i < loops; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//then
	for i := 0; i < loops; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance of th transaction

	updatedFromAccount, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	// account 1 is fromAccount and account 2 is toAccount
	require.Equal(t, account1.Balance, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance, updatedToAccount.Balance)
}
