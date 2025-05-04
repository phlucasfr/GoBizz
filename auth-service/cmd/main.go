package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/infra/cache"
	"auth-service/internal/infra/database"
	"auth-service/internal/infra/grpc/links"
	"auth-service/internal/infra/repository"
	"auth-service/internal/infra/server"
	"auth-service/internal/logger"
	"auth-service/utils"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	logger.Initialize(os.Getenv("ENVIRONMENT"))
	loadEnvironment()
}

func loadEnvironment() {
	if os.Getenv("ENVIRONMENT") == "production" {
		return
	}

	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("Failed to load .env file",
			zap.Error(err),
			zap.String("component", "config"),
		)
	}

	utils.LoadEnvInstance()
}

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

func main() {
	logger.Log.Info("Starting auth service...")

	db, err := initPostgres()
	if err != nil {
		logger.Log.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer db.Close()
	logger.Log.Info("Successfully connected to PostgreSQL")

	rdb, err := initRedis()
	if err != nil {
		logger.Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Log.Info("Successfully connected to Redis")

	linksClientRead, err := links.NewClient()
	if err != nil {
		logger.Log.Fatal("Failed to connect to links read service", zap.Error(err))
	}
	defer linksClientRead.CloseRead()
	logger.Log.Info("Successfully connected to links read service")

	linksClientWrite, err := links.NewClient()
	if err != nil {
		logger.Log.Fatal("Failed to connect to links write service", zap.Error(err))
	}
	defer linksClientWrite.CloseWrite()
	logger.Log.Info("Successfully connected to links write service")

	customerRepo := repository.NewCustomerRepository(db, rdb)
	customerHandler := handlers.NewCustomerHandler(customerRepo)

	linksHandler := handlers.NewLinksHandler(linksClientWrite, linksClientRead)

	app := server.InitFiber(customerHandler, linksHandler, rdb)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Log.Info("Server starting on port :3000")
		if err := app.Listen(":3000"); err != nil {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rdb.Close(); err != nil {
		logger.Log.Error("Error closing Redis connection", zap.Error(err))
	}

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server gracefully stopped")
}
