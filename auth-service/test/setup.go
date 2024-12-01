package test

import (
	"auth-service/internal/infra/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestRepositories struct {
	*repository.CompanyRepository
	*repository.SessionRepository
}

func SetupTestContainers(t *testing.T) (*TestRepositories, error) {
	ctx := context.Background()

	redisContainer, err := SetupRedisTestContainer(ctx)
	require.NoError(t, err)

	postgresContainer, err := SetupPostgresTestContainer(ctx)
	require.NoError(t, err)

	sessionRepo := repository.NewSessionRepository(redisContainer.Client)
	companyRepo := repository.NewCompanyRepository(postgresContainer.Pool, redisContainer.Client)

	return &TestRepositories{
		CompanyRepository: companyRepo,
		SessionRepository: sessionRepo,
	}, nil
}
