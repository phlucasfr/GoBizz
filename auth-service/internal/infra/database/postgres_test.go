package database

import (
	"auth-service/utils"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	err := godotenv.Load("../../../.env")
	require.NoError(t, err, "Failed to load .env file")
	utils.LoadEnvInstance()

	log.Println("DB_SOURCE:", utils.ConfigInstance.DBSource)

	db, err := NewPostgresConnection()
	require.NoError(t, err, "Unexpected error while initializing the Postgres connection")
	require.NotEmpty(t, db, "Database connection should not be empty")
}
