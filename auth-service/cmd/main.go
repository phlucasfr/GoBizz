package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/infra/cache"
	"auth-service/internal/infra/database"
	"auth-service/internal/infra/grpc/links"
	"auth-service/internal/infra/repository"
	"auth-service/internal/infra/server"
	"auth-service/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func initPostgres() (*pgxpool.Pool, error) {
	db, err := database.NewPostgresConnection()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initRedis() (*redis.Client, error) {
	rdb, err := cache.NewRedisClient(utils.ConfigInstance.RedisHost, utils.ConfigInstance.RedisPort)
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting current directory: %v", err)
		}

		envPath := filepath.Join(currentDir, ".env")
		log.Printf("Attempting to load .env from: %s", envPath)

		err = godotenv.Load(envPath)
		if err != nil {
			log.Fatalf("error during loading .env: %v", err)
		}
	}

	utils.LoadEnvInstance()
}

func main() {
	logger := log.New(os.Stdout, "[AUTH-SERVICE] ", log.LstdFlags|log.Lshortfile)

	logger.Println("Starting auth service...")

	db, err := initPostgres()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Println("Successfully connected to PostgreSQL")

	rdb, err := initRedis()
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	logger.Println("Successfully connected to Redis")

	linksClient, err := links.NewClient(utils.ConfigInstance.LinksServiceUrl)
	if err != nil {
		logger.Fatalf("Failed to connect to links service: %v", err)
	}
	defer linksClient.Close()
	logger.Println("Successfully connected to links service")

	customerRepo := repository.NewCustomerRepository(db, rdb)
	customerHandler := handlers.NewCustomerHandler(customerRepo)

	linksHandler := handlers.NewLinksHandler(linksClient)

	app := server.InitFiber(customerHandler, linksHandler, rdb)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Println("Server starting on port :3000")
		if err := app.Listen(":3000"); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rdb.Close(); err != nil {
		logger.Printf("Error closing Redis connection: %v", err)
	}

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server gracefully stopped")
}
