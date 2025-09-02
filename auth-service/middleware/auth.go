package middleware

import (
	"auth-service/internal/logger"
	"auth-service/utils"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a middleware function for the Fiber framework that handles
// authentication and token validation. It performs the following tasks:
func AuthMiddleware(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Log.Error("Authorization header is missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Log.Error("Invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		tokenKey := fmt.Sprintf("blacklisted-token:%s", token)
		exists, err := rdb.Exists(c.Context(), tokenKey).Result()
		if err != nil {
			logger.Log.Error("Error checking token blacklist")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		if exists == 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has been invalidated",
			})
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			logger.Log.Error("Invalid or expired token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		logger.Log.Info("User authenticated successfully")
		return c.Next()
	}
}
