package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func GenerateSecureToken(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid length")
	}

	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error during generating safe token: %w", err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func GenerateSessionToken() string {
	token, err := GenerateSecureToken(32)
	if err != nil {
		panic("critical failure during generating safe token: " + err.Error())
	}
	return token
}
