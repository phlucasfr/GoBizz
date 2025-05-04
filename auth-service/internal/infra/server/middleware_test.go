package server

import (
	"auth-service/internal/logger"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func TestJWTMiddleware(t *testing.T) {
	app := fiber.New()
	masterKey := "01234567890123456789012345678901"

	app.Use(JWTMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	t.Run("Valid Token", func(t *testing.T) {
		token, err := utils.GenerateJWT("123", masterKey)
		require.NoError(t, err, "Expected no error generating the token")
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", token)

		resp, err := app.Test(req)
		require.NoError(t, err, "Expected no error executing the request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status OK for a valid token")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "invalid_token")

		resp, err := app.Test(req)
		require.NoError(t, err, "Expected no error executing the request")
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status Unauthorized for an invalid token")
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Create an expired token manually using jwt.RegisteredClaims
		claims := jwt.RegisteredClaims{
			Subject:   "123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		}
		tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredToken, err := tokenObj.SignedString([]byte(masterKey))
		require.NoError(t, err, "Expected no error generating an expired token")

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", expiredToken)

		resp, err := app.Test(req)
		require.NoError(t, err, "Expected no error executing the request")
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status Unauthorized for an expired token")
	})
}
