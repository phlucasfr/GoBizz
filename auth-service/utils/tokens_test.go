package utils

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestGenerateSecureToken(t *testing.T) {
	t.Run("Generates token with correct length", func(t *testing.T) {
		token, err := GenerateSecureToken(32)
		require.NoError(t, err, "Expected no error when generating token")
		require.Equal(t, 43, utf8.RuneCountInString(token), "Expected token length to be 43")
	})

	t.Run("Tokens are unique", func(t *testing.T) {
		t1, err1 := GenerateSecureToken(32)
		require.NoError(t, err1, "Expected no error when generating first token")
		t2, err2 := GenerateSecureToken(32)
		require.NoError(t, err2, "Expected no error when generating second token")
		require.NotEqual(t, t1, t2, "Tokens should be unique")
	})

	t.Run("Error with invalid size", func(t *testing.T) {
		_, err := GenerateSecureToken(-1)
		require.Error(t, err, "Expected an error when generating token with invalid size")
	})
}
