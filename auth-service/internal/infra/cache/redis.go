package cache

import (
	"auth-service/pkg/util"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() (*redis.Client, error) {
	env := util.GetConfig(".")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Redis: %v", err)
	}
	return rdb, nil
}
