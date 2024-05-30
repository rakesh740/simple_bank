package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := CreateTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
		Status:        util.RandomCurrency(),
	}

	Transfer, err := testQueries.CreateTransfers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	require.Equal(t, args.FromAccountID, Transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, Transfer.ToAccountID)
	require.Equal(t, args.Amount, Transfer.Amount)
	require.Equal(t, args.Status, Transfer.Status)
	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}

func Test_CreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func Test_GetTransfer(t *testing.T) {

	argTransfer := createRandomTransfer(t)

	Transfer, err := testQueries.GetTransfer(context.Background(), argTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	require.Equal(t, argTransfer.FromAccountID, Transfer.FromAccountID)
	require.Equal(t, argTransfer.ToAccountID, Transfer.ToAccountID)
	require.Equal(t, argTransfer.Amount, Transfer.Amount)
	require.Equal(t, argTransfer.Status, Transfer.Status)
	require.WithinDuration(t, argTransfer.CreatedAt, Transfer.CreatedAt, time.Second)
	require.Equal(t, argTransfer.ID, Transfer.ID)
}

func Test_ListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}
	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, a := range accounts {
		require.NotEmpty(t, a)
	}
}
