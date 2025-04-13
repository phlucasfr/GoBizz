package middleware

import (
	"auth-service/utils"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header is missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		tokenKey := fmt.Sprintf("blacklisted-token:%s", token)
		exists, err := rdb.Exists(c.Context(), tokenKey).Result()
		if err != nil {
			// Handle Redis error
			log.Printf("Error checking token blacklist: %v", err)
		} else if exists == 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has been invalidated",
			})
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			log.Printf("Error validating token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
