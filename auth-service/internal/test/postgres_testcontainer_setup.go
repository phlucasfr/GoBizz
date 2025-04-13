package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "test_db"
	DbUser = "test_user"
	DbPass = "test_password"
)

type TestPostgres struct {
	Pool      *pgxpool.Pool
	DbAddress string
	container testcontainers.Container
}

func SetupPostgresTestContainer(ctx context.Context) (*TestPostgres, error) {
	container, pool, dbAddr, err := createContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when configuring the test: %w", err)
	}

	err = migrateDb(dbAddr)
	if err != nil {
		return nil, fmt.Errorf("error when running migration: %w", err)
	}

	return &TestPostgres{
		Pool:      pool,
		DbAddress: dbAddr,
		container: container,
	}, nil
}

func (tdb *TestPostgres) TearDown() error {
	tdb.Pool.Close()
	return tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, nil, "", fmt.Errorf("error when starting the container: %w", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, "", fmt.Errorf("error when getting the container's external port: %w", err)
	}

	log.Println("PostgreSQL container is ready on port:", p.Port())

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())

	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName))
	if err != nil {
		return nil, nil, dbAddr, fmt.Errorf("error parsing pgxpool config: %w", err)
	}

	var pool *pgxpool.Pool
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			err = pool.Ping(ctx)
		}

		if err == nil {
			break
		}

		log.Println("Waiting for database connection...")
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, nil, dbAddr, fmt.Errorf("error connecting to database: %w", err)
	}

	return container, pool, dbAddr, nil
}

func migrateDb(dbAddr string) error {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get file path")
	}
	pathToMigrationFiles := filepath.Join(filepath.Dir(path), "../../migrations")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
