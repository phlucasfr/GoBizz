package utils

import (
	"log"
	"os"
)

type Config struct {
	DBSource        string
	MasterKey       string
	RedisPort       string
	RedisHost       string
	AllowedOrigins  string
	FrontendSource  string
	SendGridApiKey  string
	LinksServiceUrl string
}

var (
	ConfigInstance Config
)

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadEnvInstance() {
	log.Println("Loading environment variables...")

	ConfigInstance = Config{
		DBSource:        getEnvWithDefault("DB_SOURCE", "postgresql://postgres:postgres@localhost:5432/gobizz?sslmode=disable"),
		MasterKey:       getEnvWithDefault("MASTER_KEY", "j8Lr5a2W9cX3pZb7Nq4eK6dF1gHtR0mU"),
		RedisPort:       getEnvWithDefault("REDIS_PORT", "6379"),
		RedisHost:       getEnvWithDefault("REDIS_HOST", "localhost"),
		AllowedOrigins:  getEnvWithDefault("ALLOWED_ORIGINS", "http://localhost,http://127.0.0.1,http://localhost:5173,http://localhost:3000"),
		FrontendSource:  getEnvWithDefault("FRONTEND_SOURCE", "http://localhost:5173"),
		SendGridApiKey:  os.Getenv("SENDGRID_API_KEY"),
		LinksServiceUrl: getEnvWithDefault("LINKS_SERVICE_URL", "localhost:50051"),
	}

	log.Printf("Configuration loaded successfully:")
	log.Printf("Redis configuration - Host: %s, Port: %s",
		ConfigInstance.RedisHost,
		ConfigInstance.RedisPort)
}
