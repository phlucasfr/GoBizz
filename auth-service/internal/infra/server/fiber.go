package server

import (
	"auth-service/internal/handlers"
	"auth-service/utils"
	"os"

	json "github.com/goccy/go-json"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// InitFiber initializes and configures a new Fiber application instance.
// It sets up middleware for CORS, logging, error recovery, and encryption,
// and registers application routes.
//
// Parameters:
//   - customerHandler: A pointer to the CustomerHandler, responsible for handling customer-related routes.
//   - linksHandler: A pointer to the LinksHandler, responsible for handling link-related routes.
//   - rdb: A pointer to a Redis client instance for caching or other Redis-related operations.
//
// Returns:
//   - *fiber.App: A fully configured Fiber application instance ready to start serving requests.
func InitFiber(customerHandler *handlers.CustomerHandler, linksHandler *handlers.LinksHandler, eventsHandler *handlers.EventsHandler, rdb *redis.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:     "auth-service API",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     utils.ConfigInstance.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
	}))

	app.Use(recover.New())

	if os.Getenv("ENVIRONMENT") == "production" {
		app.Use(func(c *fiber.Ctx) error {
			if c.Protocol() == "http" {
				return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusPermanentRedirect)
			}
			return c.Next()
		})

		app.Use(func(c *fiber.Ctx) error {
			c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			return c.Next()
		})
	}

	app.Use(logger.New())
	app.Use(EncryptionMiddleware())

	setupRoutes(app, customerHandler, linksHandler, eventsHandler, rdb)
	return app
}
