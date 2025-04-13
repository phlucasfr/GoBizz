package server

import (
	"auth-service/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func EncryptionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > 0 {
			if decryptErr := decryptRequest(c, utils.ConfigInstance.MasterKey); decryptErr != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Failed to decrypt request: " + decryptErr.Error(),
				})
			}
		}

		err := c.Next()

		if err == nil && c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			if encryptErr := encryptResponse(c, utils.ConfigInstance.MasterKey); encryptErr != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to encrypt response",
				})
			}
		}

		return err
	}
}

func decryptRequest(c *fiber.Ctx, key string) error {
	if len(c.Body()) == 0 {
		return nil
	}

	var payload struct {
		Data string `json:"data"`
	}

	if err := json.Unmarshal(c.Body(), &payload); err == nil && payload.Data != "" {
		encrypted := payload.Data
		decrypted, err := utils.Decrypt(encrypted, key)
		if err != nil {
			return fmt.Errorf("erro na descriptografia: %w", err)
		}
		c.Request().SetBody([]byte(decrypted))
		return nil
	}

	encrypted := string(c.Body())
	decrypted, err := utils.Decrypt(encrypted, key)
	if err != nil {
		return fmt.Errorf("erro na descriptografia: %w", err)
	}

	c.Request().SetBody([]byte(decrypted))
	return nil
}

func encryptResponse(c *fiber.Ctx, key string) error {
	originalBody := c.Response().Body()
	contentType := c.Response().Header.ContentType()

	if len(originalBody) == 0 || string(contentType) != fiber.MIMEApplicationJSON {
		return nil
	}

	var jsonCheck interface{}
	if err := json.Unmarshal(originalBody, &jsonCheck); err != nil {
		return fmt.Errorf("invalid JSON response: %v", err)
	}

	encrypted, err := utils.Encrypt(string(originalBody), key)
	if err != nil {
		return err
	}

	response := struct{ Data string }{Data: encrypted}
	jsonResponse, _ := json.Marshal(response)

	c.Response().SetBody(jsonResponse)
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
	return nil
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token not provided",
			})
		}

		userID, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
