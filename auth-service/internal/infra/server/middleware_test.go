package server

import (
	"auth-service/internal/logger"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type testResponse struct {
	Message string `json:"message"`
}

func TestMain(m *testing.M) {
	logger.Initialize("development")
	code := m.Run()
	logger.Sync()
	os.Exit(code)
}

func TestEncryptionMiddleware(t *testing.T) {
	app := fiber.New()
	masterKey := utils.ConfigInstance.MasterKey

	app.Use(EncryptionMiddleware())
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.JSON(testResponse{Message: "success"})
	})

	t.Run("Should encrypt/decrypt complete flow", func(t *testing.T) {
		originalPayload := "secret message"

		encrypted, err := utils.Encrypt(originalPayload, masterKey)
		require.NoError(t, err, "Encryption should succeed")

		reqBody := struct{ Data string }{Data: encrypted}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		var response struct{ Data string }
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err, "Response decoding should succeed")

		decryptedResponse, err := utils.Decrypt(response.Data, masterKey)
		require.NoError(t, err, "Decryption should succeed")

		var result testResponse
		err = json.Unmarshal([]byte(decryptedResponse), &result)
		require.NoError(t, err, "Unmarshal should succeed")
		require.Equal(t, "success", result.Message, "Response message should match")
	})

	t.Run("Should handle invalid encrypted payload", func(t *testing.T) {
		invalidPayload := []byte(`{"data":"invalid_encrypted_string"}`)
		req := httptest.NewRequest("POST", "/test", bytes.NewReader(invalidPayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return bad request status")
	})

	t.Run("Should handle empty body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Should handle empty body")
	})
}
