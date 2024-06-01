package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	user1 := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomPassword(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomName(),
	}

	user, err := testQueries.CreateUser(context.Background(), user1)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user1.Username, user.Username)
	require.Equal(t, user1.HashedPassword, user.HashedPassword)
	require.Equal(t, user1.FullName, user.FullName)
	require.Equal(t, user1.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user

}
func Test_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func Test_GetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user1.Username, user.Username)
	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, user1.FullName, user.FullName)
	require.Equal(t, user1.HashedPassword, user.HashedPassword)
	require.WithinDuration(t, user1.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user.PasswordChangedAt, time.Second)

}
