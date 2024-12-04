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
	"auth-service/pkg/util"
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

func initServer(companyHandler *handlers.CompanyHandler, sessionHandler *handlers.SessionHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Your Company Management API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	env := util.GetConfig(".")

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowOrigins:     env.AllowedOrigins,
		AllowCredentials: true,
	}))

	v1 := app.Group("/v1")
	v1.Get("/companies/:id", companyHandler.GetByID)

	v1.Put("/companies/reset-password", companyHandler.ResetPassword)
	v1.Put("/companies/email-verification", companyHandler.VerifyCompanyByEmail)

	v1.Post("/companies", companyHandler.Create)
	v1.Post("/companies/login", companyHandler.Login)
	v1.Post("/companies/recovery", companyHandler.RecoverPassword)
	v1.Post("/companies/email-verification", companyHandler.SendVerificationEmail)

	v1.Get("/sessions", sessionHandler.ValidateSession)

	v1.Post("/sessions", sessionHandler.CreateSession)

	v1.Delete("/sessions", sessionHandler.DeleteSession)

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

	sessionRepo := repository.NewSessionRepository(rdb)
	companyRepo := repository.NewCompanyRepository(db, rdb)

	sessionHandler := handlers.NewSessionHandler(sessionRepo)
	companyHandler := handlers.NewCompanyHandler(companyRepo, sessionRepo)

	app := initServer(companyHandler, sessionHandler)

	log.Fatal(app.Listen(":3000"))
}
