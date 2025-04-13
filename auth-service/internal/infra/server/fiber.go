package server

import (
	"auth-service/internal/handlers"
	"auth-service/utils"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func InitFiber(customerHandler *handlers.CustomerHandler, linksHandler *handlers.LinksHandler, rdb *redis.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "auth-service API",
		Prefork: true,
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

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(EncryptionMiddleware())

	setupRoutes(app, customerHandler, linksHandler, rdb)
	return app
}
