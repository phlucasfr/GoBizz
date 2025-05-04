package server

import (
	"auth-service/internal/logger"
	"auth-service/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// EncryptionMiddleware is a middleware for the Fiber framework that handles
// encryption and decryption of HTTP requests and responses. It performs the
// following tasks:
//
//  1. Decrypts the incoming request body using a master key before passing it
//     to the next handler in the middleware chain. If decryption fails, it
//     responds with a 400 Bad Request status and an error message.
//
//  2. After the next handler processes the request, it checks if the response
//     status code indicates success (2xx). If so, it encrypts the response body
//     using the same master key. If encryption fails, it responds with a 500
//     Internal Server Error status and an error message.
//
// This middleware ensures that sensitive data is securely handled during
// transmission. It relies on utility functions for encryption and decryption
// and uses a master key from the configuration instance.
func EncryptionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > 0 {
			if decryptErr := decryptRequest(c, utils.ConfigInstance.MasterKey); decryptErr != nil {
				logger.Log.Error("Failed to decrypt request", zap.Error(decryptErr))
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Failed to decrypt request: " + decryptErr.Error(),
				})
			}
		}

		err := c.Next()

		if err == nil && c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			if encryptErr := encryptResponse(c, utils.ConfigInstance.MasterKey); encryptErr != nil {
				logger.Log.Error("Failed to encrypt response", zap.Error(encryptErr))
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to encrypt response",
				})
			}
		}

		logger.Log.Info("Request processed successfully")
		return err
	}
}

func decryptRequest(c *fiber.Ctx, key string) error {
	if len(c.Body()) == 0 {
		logger.Log.Info("Request body is empty, skipping decryption")
		return nil
	}

	var payload struct {
		Data string `json:"data"`
	}

	if err := json.Unmarshal(c.Body(), &payload); err == nil && payload.Data != "" {
		encrypted := payload.Data
		decrypted, err := utils.Decrypt(encrypted, key)
		if err != nil {
			logger.Log.Error("Failed to decrypt request body", zap.Error(err))
			return fmt.Errorf("decryption error: %w", err)
		}
		c.Request().SetBody([]byte(decrypted))
		logger.Log.Info("Request body decrypted successfully")
		return nil
	}

	encrypted := string(c.Body())
	decrypted, err := utils.Decrypt(encrypted, key)
	if err != nil {
		logger.Log.Error("Failed to decrypt request body", zap.Error(err))
		return fmt.Errorf("decryption error: %w", err)
	}

	c.Request().SetBody([]byte(decrypted))
	logger.Log.Info("Request body decrypted successfully")
	return nil
}

func encryptResponse(c *fiber.Ctx, key string) error {
	originalBody := c.Response().Body()
	contentType := c.Response().Header.ContentType()

	if len(originalBody) == 0 || string(contentType) != fiber.MIMEApplicationJSON {
		logger.Log.Info("Response body is empty or not JSON, skipping encryption")
		return nil
	}

	var jsonCheck interface{}
	if err := json.Unmarshal(originalBody, &jsonCheck); err != nil {
		logger.Log.Error("Failed to unmarshal JSON response", zap.Error(err))
		return fmt.Errorf("invalid JSON response: %v", err)
	}

	encrypted, err := utils.Encrypt(string(originalBody), key)
	if err != nil {
		logger.Log.Error("Failed to encrypt response body", zap.Error(err))
		return err
	}

	response := struct{ Data string }{Data: encrypted}
	jsonResponse, _ := json.Marshal(response)

	c.Response().SetBody(jsonResponse)
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)

	logger.Log.Info("Response body encrypted successfully")
	return nil
}

// JWTMiddleware is a middleware function for the Fiber framework that validates
// JSON Web Tokens (JWT) from the "Authorization" header of incoming requests.
//
// If the token is not provided or is invalid, the middleware responds with a
// 401 Unauthorized status and an appropriate error message in JSON format.
//
// Upon successful validation, the middleware extracts the user ID from the token
// and stores it in the request's local context under the key "userID", allowing
// subsequent handlers to access it.
//
// Returns:
//   - A Fiber handler function to be used as middleware.
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			logger.Log.Error("Authorization token not provided")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token not provided",
			})
		}

		userID, err := utils.ValidateJWT(token)
		if err != nil {
			logger.Log.Error("Invalid token", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", userID)
		logger.Log.Info("Token validated successfully", zap.String("userID", userID.UserID))
		return c.Next()
	}
}
