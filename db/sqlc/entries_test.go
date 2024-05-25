package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {

	account := createRandomAccount(t)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
		Status:    util.RandomCurrency(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.Equal(t, args.Status, entry.Status)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func Test_CreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func Test_GetEntry(t *testing.T) {

	argEntry := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), argEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, argEntry.AccountID, entry.AccountID)
	require.Equal(t, argEntry.Amount, entry.Amount)
	require.Equal(t, argEntry.Status, entry.Status)
	require.WithinDuration(t, argEntry.CreatedAt, entry.CreatedAt, time.Second)
	require.Equal(t, argEntry.ID, entry.ID)
}


func Test_ListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}
	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, a := range accounts {
		require.NotEmpty(t, a)
	}
}