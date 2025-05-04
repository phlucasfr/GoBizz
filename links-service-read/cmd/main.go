package main

import (
	"context"
	"fmt"
	"links-service-read/internal/infra/database"
	"links-service-read/internal/infra/repository"
	"links-service-read/internal/logger"
	"links-service-read/internal/server"
	"links-service-read/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func initDynamo() (*dynamodb.Client, error) {
	client, err := database.NewDynamoClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create DynamoDB client: %v", err)
	}

	return client, nil
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := initDynamo()
	if err != nil {
		logger.Log.Fatal("Failed to connect to DynamoDB",
			zap.Error(err),
			zap.String("component", "database"),
		)
	}

	logger.Log.Info("Successfully connected to DynamoDB",
		zap.String("component", "database"),
	)

	linksRepo := repository.NewLinksRepository(db)
	logger.Log.Info("LinksRepository initialized",
		zap.String("component", "repository"),
	)

	go func() {
		logger.Log.Info("Starting gRPC server",
			zap.String("port", "50051"),
			zap.String("component", "server"),
		)
		if err := server.StartGRPCServer("50051", linksRepo); err != nil {
			logger.Log.Error("Failed to start gRPC server",
				zap.Error(err),
				zap.String("component", "server"),
			)
			cancel()
		}
	}()

	select {
	case sig := <-sigChan:
		logger.Log.Info("Received shutdown signal",
			zap.String("signal", sig.String()),
			zap.String("component", "lifecycle"),
		)
	case <-ctx.Done():
		logger.Log.Warn("Context cancelled",
			zap.String("component", "lifecycle"),
		)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	<-shutdownCtx.Done()
	logger.Log.Info("Server shutdown complete",
		zap.String("component", "lifecycle"),
	)
}
