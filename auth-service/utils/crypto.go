package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// Encrypt encrypts data using NaCl secretbox
func Encrypt(data, key string) (string, error) {

	if data == "" {
		return "", errors.New("data is empty")
	}

	// Gera o nonce (24 bytes)
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", errors.New("error generating nonce")
	}

	// Converte a chave para 32 bytes
	var keyArr [32]byte
	copy(keyArr[:], key)

	// Criptografa os dados (incluindo a tag de autenticação)
	encrypted := secretbox.Seal(nil, []byte(data), &nonce, &keyArr)

	// Concatena nonce + ciphertext + tag
	combined := make([]byte, 24+len(encrypted))
	copy(combined[:24], nonce[:])  // Primeiros 24 bytes: nonce
	copy(combined[24:], encrypted) // Restante: ciphertext + tag (16 bytes)

	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decrypts data encrypted with NaCl secretbox
func Decrypt(encrypted, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < 24 {
		return "", errors.New("ciphertext too short")
	}

	var nonce [24]byte
	copy(nonce[:], data[:24])

	var keyArr [32]byte
	copy(keyArr[:], key)

	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, &keyArr)
	if !ok {
		return "", errors.New("decryption failed")
	}

	return string(decrypted), nil
}
