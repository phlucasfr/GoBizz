package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedisTestContainer(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := SetupRedisTestContainer(ctx)
	require.NoError(t, err, "failed to set up Redis test container")

	defer func() {
		err := redisContainer.TearDown()
		require.NoError(t, err, "failed to tear down Redis test container")
	}()
}

func TestPostgresTestContainer(t *testing.T) {
	ctx := context.Background()

	testDB, err := SetupPostgresTestContainer(ctx)
	require.NoError(t, err)

	defer func() {
		err = testDB.TearDown()
		require.NoError(t, err)
	}()
}
