package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func Test_Password(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "correct password",
			password: RandomPassword(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)
			require.NoError(t, err)

			err = ComparePassword(tt.password, hashedPassword)
			require.NoError(t, err)
		})
	}
}

func Test_IncorrectPassword(t *testing.T) {

	password := RandomPassword()
	incorrectPassword := RandomPassword()

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = ComparePassword(incorrectPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
