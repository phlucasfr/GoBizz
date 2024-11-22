package main

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-service/internal/handlers"
	"auth-service/internal/infra/cache"
	"auth-service/internal/infra/database"
	"auth-service/internal/infra/repository"
)

func initDB() (*pgxpool.Pool, error) {
	db, err := database.NewPostgresConnection()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initRedis() (*redis.Client, error) {
	rdb, err := cache.NewRedisClient()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func initServer(companyHandler *handlers.CompanyHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Your Company Management API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	v1 := app.Group("/v1")
	v1.Post("/companies", companyHandler.Create)
	v1.Post("/companies/sms/verify", companyHandler.VerifyCompanyBySms)
	v1.Get("/companies/:id", companyHandler.GetByID)

	return app
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rdb, err := initRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	companyRepo := repository.NewCompanyRepository(db, rdb)

	companyHandler := handlers.NewCompanyHandler(companyRepo)

	app := initServer(companyHandler)

	log.Fatal(app.Listen(":3000"))
}
