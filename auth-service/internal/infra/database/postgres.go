package database

import (
	"auth-service/internal/logger"
	"auth-service/utils"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewPostgresConnection establishes a new connection pool to a PostgreSQL database
// using the pgxpool package. It parses the database configuration, sets connection
// pool parameters, and ensures the connection is valid by pinging the database.
//
// Returns:
//   - *pgxpool.Pool: A pointer to the connection pool if the connection is successful.
//   - error: An error if the connection could not be established or the database could not be pinged.
//
// Connection Pool Configuration:
//   - MaxConns: Maximum number of connections in the pool (50).
//   - MinConns: Minimum number of connections in the pool (10).
//   - MaxConnLifetime: Maximum lifetime of a connection (10 minutes).
//   - MaxConnIdleTime: Maximum idle time for a connection (5 minutes).
//   - HealthCheckPeriod: Interval for health checks on idle connections (30 minutes).
//
// Context:
//
//	A timeout of 10 seconds is applied when establishing the connection pool.
//
// Errors:
//   - Returns an error if the database URL cannot be parsed.
//   - Returns an error if the connection pool cannot be created.
//   - Returns an error if the database cannot be pinged.
func NewPostgresConnection() (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(utils.ConfigInstance.DBSource)
	if err != nil {
		logger.Log.Error("Unable to parse database URL", zap.Error(err))
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Log.Error("Unable to create connection pool", zap.Error(err))
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Log.Error("Unable to ping database", zap.Error(err))
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	logger.Log.Info("PostgreSQL connection pool created successfully")
	return pool, nil
}
