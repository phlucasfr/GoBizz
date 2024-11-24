package test

import (
	"auth-service/test"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedisTestContainer(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := test.SetupRedisTestContainer(ctx)
	require.NoError(t, err)

	defer func() {
		err := redisContainer.TearDown()
		require.NoError(t, err)
	}()
}

func TestPostgresTestContainer(t *testing.T) {
	ctx := context.Background()

	testDB, err := test.SetupPostgresTestContainer(ctx)
	require.NoError(t, err)

	defer func() {
		err = testDB.TearDown()
		require.NoError(t, err)
	}()
}
