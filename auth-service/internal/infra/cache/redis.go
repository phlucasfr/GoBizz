package cache

import (
	"auth-service/pkg/util"
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func getRedisPassword() string {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "production" {
		return os.Getenv("REDIS_PASSWORD")
	}
	return ""
}

func NewRedisClient() (*redis.Client, error) {
	env := util.GetConfig(".")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Password: getRedisPassword(),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Redis: %v", err)
	}
	return rdb, nil
}
