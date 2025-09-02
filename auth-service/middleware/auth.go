package middleware

import (
	"auth-service/internal/logger"
	"auth-service/utils"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func clientIP(c *fiber.Ctx) string {
	if ip := strings.TrimSpace(c.Get("CF-Connecting-IP")); ip != "" {
		return normalizeIP(ip)
	}

	if xff := strings.TrimSpace(c.Get(fiber.HeaderXForwardedFor)); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return normalizeIP(strings.TrimSpace(parts[0]))
		}
	}

	return normalizeIP(c.IP())
}

func clientUA(c *fiber.Ctx) string {
	return strings.ToLower(strings.TrimSpace(c.Get("User-Agent")))
}

func normalizeIP(ip string) string {
	// Normalize IPv6-mapped IPv4: ::ffff:1.2.3.4 -> 1.2.3.4
	if strings.HasPrefix(ip, "::ffff:") && len(ip) > len("::ffff:") {
		return strings.TrimPrefix(ip, "::ffff:")
	}
	return ip
}

// AuthMiddleware is a middleware function for the Fiber framework that handles
// authentication and token validation. It performs the following tasks:
//
//  1. Checks for the presence of the "Authorization" header in the incoming request.
//     If the header is missing, it responds with a 401 Unauthorized status and an error message.
//
//  2. Validates the format of the "Authorization" header to ensure it starts with "Bearer ".
//     If the format is invalid, it responds with a 401 Unauthorized status and an error message.
//
//  3. Extracts the token from the "Authorization" header and checks if the token is blacklisted
//     in the Redis database. If the token is blacklisted, it responds with a 401 Unauthorized status
//     and an error message.
//
//  4. Validates the token using the `utils.ValidateJWT` function. If the token is invalid or expired,
//     it responds with a 401 Unauthorized status and an error message.
//
//  5. If the token is valid, it extracts user information (e.g., user ID and email) from the token
//     claims and stores them in the request context using `c.Locals`.
//
// 6. Proceeds to the next middleware or handler in the chain if authentication is successful.
//
// Parameters:
// - rdb (*redis.Client): A Redis client instance used to check for blacklisted tokens.
//
// Returns:
// - fiber.Handler: A Fiber middleware handler function.
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
		} else if exists == 1 {
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

		ip := clientIP(c)
		ua := clientUA(c)

		calculated := utils.ComputeIpHash(ip, ua)
		if calculated != claims.IpHash {
			logger.Log.Warn("IP/User-Agent mismatch on token",
				zap.String("ip", ip),
				zap.String("ua", ua),
				zap.String("calc_hash", calculated),
				zap.String("claim_hash", claims.IpHash),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token context",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		logger.Log.Info("User authenticated successfully")
		return c.Next()
	}
}
