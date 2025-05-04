package utils

import "crypto/rand"

// GenerateRandomSlug generates a random alphanumeric string of the specified length.
// The generated string is composed of characters from the charset "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".
// It uses a cryptographically secure random number generator to ensure randomness.
//
// Parameters:
//   - length: The desired length of the generated slug.
//
// Returns:
//   - A randomly generated alphanumeric string of the specified length.
//   - An error if the random number generator fails.
func GenerateRandomSlug(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
