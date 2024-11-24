package test

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestRedis struct {
	Container testcontainers.Container
	Client    *redis.Client
}

func SetupRedisTestContainer(ctx context.Context) (*TestRedis, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("error when starting Redis container: %w", err)
	}

	log.Println("Redis container started successfully")

	host, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Redis container host: %w", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return nil, fmt.Errorf("error getting Redis container port: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	return &TestRedis{Container: redisContainer, Client: rdb}, nil
}

func (tr *TestRedis) TearDown() error {
	ctx := context.Background()
	if err := tr.Client.Close(); err != nil {
		return err
	}
	return tr.Container.Terminate(ctx)
}
