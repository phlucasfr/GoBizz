package main

import (
	"context"
	"links-service/internal/infra/database"
	"links-service/internal/infra/repository"
	"links-service/internal/server"
	"links-service/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("error during loading .env: %v", err)
	}
	utils.LoadEnvInstance()
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := initPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to PostgreSQL")

	repo := repository.New(db)

	go func() {
		log.Println("Starting gRPC server on port 50051...")
		if err := server.StartGRPCServer("50051", repo); err != nil {
			log.Printf("Failed to start gRPC server: %v", err)
			cancel()
		}
	}()

	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	case <-ctx.Done():
		log.Println("Context canceled, shutting down...")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	<-shutdownCtx.Done()
	log.Println("Server shutdown complete")
}
