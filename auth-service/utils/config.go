package utils

import (
	"log"
	"os"
)

type Config struct {
	DBSource             string
	MasterKey            string
	RedisPort            string
	RedisHost            string
	AllowedOrigins       string
	FrontendSource       string
	SendGridApiKey       string
	LinksServiceReadUrl  string
	LinksServiceWriteUrl string
}

var (
	ConfigInstance Config
)

func LoadEnvInstance() {
	log.Println("Loading environment variables...")

	ConfigInstance = Config{
		DBSource:             os.Getenv("DB_SOURCE"),
		MasterKey:            os.Getenv("MASTER_KEY"),
		RedisPort:            os.Getenv("REDIS_PORT"),
		RedisHost:            os.Getenv("REDIS_HOST"),
		AllowedOrigins:       os.Getenv("ALLOWED_ORIGINS"),
		FrontendSource:       os.Getenv("FRONTEND_SOURCE"),
		SendGridApiKey:       os.Getenv("SENDGRID_API_KEY"),
		LinksServiceReadUrl:  os.Getenv("LINKS_SERVICE_READ_URL"),
		LinksServiceWriteUrl: os.Getenv("LINKS_SERVICE_WRITE_URL"),
	}

	log.Printf("Configuration loaded successfully:")
	log.Printf("Redis configuration - Host: %s, Port: %s",
		ConfigInstance.RedisHost,
		ConfigInstance.RedisPort)
}
