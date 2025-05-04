package cache

import (
	"auth-service/internal/logger"
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func getRedisPassword() string {
	if os.Getenv("ENVIRONMENT") == "production" {
		return os.Getenv("REDIS_PASSWORD")
	}
	return ""
}

// NewRedisClient creates and returns a new Redis client instance.
// It attempts to connect to a Redis server using the provided host and port.
// If the host or port is empty, it returns an error indicating invalid configuration.
// The function pings the Redis server to verify the connection and returns an error
// if the connection attempt fails.
//
// Parameters:
//   - host: The hostname or IP address of the Redis server.
//   - port: The port number on which the Redis server is running.
//
// Returns:
//   - *redis.Client: A pointer to the initialized Redis client.
//   - error: An error if the connection to Redis fails or if the configuration is invalid.
func NewRedisClient(host string, port string) (*redis.Client, error) {
	logger.Log.Info("Attempting to connect to Redis", zap.String("host", host), zap.String("port", port))

	if host == "" || port == "" {
		logger.Log.Error("Invalid Redis configuration", zap.String("host", host), zap.String("port", port))
		return nil, fmt.Errorf("invalid Redis configuration: host=%s, port=%s", host, port)
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: getRedisPassword(),
	})

	logger.Log.Info("Redis client created, attempting to ping...")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis", zap.Error(err))
		return nil, fmt.Errorf("error during connect to Redis: %v", err)
	}

	logger.Log.Info("Successfully connected to Redis", zap.String("host", host), zap.String("port", port))
	return rdb, nil
}
