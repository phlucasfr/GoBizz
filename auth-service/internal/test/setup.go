package test

import (
	"auth-service/internal/infra/repository"
	"context"
	"fmt"
)

type TestRepositories struct {
	*repository.CustomerRepository
}

func SetupTestContainers() (*TestRepositories, error) {
	ctx := context.Background()

	redisContainer, err := SetupRedisTestContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("error during setup redis test container: %v", err)
	}

	postgresContainer, err := SetupPostgresTestContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("error during setup postgres test container: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(postgresContainer.Pool, redisContainer.Client)

	return &TestRepositories{
		CustomerRepository: customerRepo,
	}, nil
}
