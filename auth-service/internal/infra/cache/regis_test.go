package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
)

// TestNewRedisClient_Success starts a fake Redis server using miniredis
func TestNewRedisClient_Success(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err, "Failed to start miniredis")
	defer s.Close()

	client, err := NewRedisClient(s.Host(), s.Port())
	require.NoError(t, err, "Expected no error when connecting to miniredis")
	require.NotNil(t, client, "Redis client should not be nil")

	res, err := client.Ping(context.Background()).Result()
	require.NoError(t, err, "Ping should not return an error")
	require.Equal(t, "PONG", res, "Expected PONG response from Redis Ping")
}

// TestNewRedisClient_Failure starts a fake Redis server using miniredis
func TestNewRedisClient_Failure(t *testing.T) {
	client, err := NewRedisClient("127.0.0.1", "1234")
	require.Error(t, err, "Expected error when connecting to an unavailable Redis")
	require.Nil(t, client, "Redis client should be nil on failure")
}
