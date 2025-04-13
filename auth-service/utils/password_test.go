package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err, "Hashing password should not return an error")
	require.NotEmpty(t, hashedPassword1, "Hashed password should not be empty")

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err, "Correct password should match the hash")

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error(), "Wrong password should return a mismatched hash error")

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err, "Hashing the same password again should not return an error")
	require.NotEmpty(t, hashedPassword2, "Second hashed password should not be empty")

	require.NotEqual(t, hashedPassword1, hashedPassword2, "Hashed passwords should be different due to salting")
}
