package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "01234567890123456789012345678901"
	originalText := "secret message"

	encrypted, err := Encrypt(originalText, key)
	require.NoError(t, err, "Encryption should not return an error")
	require.NotEmpty(t, encrypted, "Encrypted text should not be empty")

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err, "Decryption should not return an error")
	require.Equal(t, originalText, decrypted, "Decrypted text should match the original")
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := "01234567890123456789012345678901"
	wrongKey := "abcdefghijklmnopqrstuvwx12345678"
	originalText := "safe message"

	encrypted, err := Encrypt(originalText, key)
	require.NoError(t, err, "Encryption should not return an error")

	_, err = Decrypt(encrypted, wrongKey)
	require.Error(t, err, "Decryption with a wrong key should return an error")
}

func TestDecryptInvalidData(t *testing.T) {
	key := "01234567890123456789012345678901"
	invalidData := "invalid_base64"

	_, err := Decrypt(invalidData, key)
	require.Error(t, err, "Decryption of invalid data should return an error")
}
