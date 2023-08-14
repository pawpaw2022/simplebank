package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// compare with the original password
	err = ComparePassword(hashedPassword1, password)
	require.NoError(t, err)

	// compare with the wrong password
	wrongPassword := RandomString(6)
	err = ComparePassword(hashedPassword1, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// rehashing should not be the same
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
