package database

import (
	"context"
	"fmt"
	"links-service/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection() (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(utils.ConfigInstance.DBSource)
	if err != nil {
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
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return pool, nil
}
