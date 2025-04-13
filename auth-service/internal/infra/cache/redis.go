package cache

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func getRedisPassword() string {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "production" {
		return os.Getenv("REDIS_PASSWORD")
	}
	return ""
}

func NewRedisClient(host string, port string) (*redis.Client, error) {
	log.Printf("Attempting to connect to Redis at %s:%s", host, port)

	if host == "" || port == "" {
		return nil, fmt.Errorf("invalid Redis configuration: host=%s, port=%s", host, port)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Redis connection address: %s", addr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: getRedisPassword(),
	})

	log.Println("Redis client created, attempting to ping...")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error during connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
	return rdb, nil
}
