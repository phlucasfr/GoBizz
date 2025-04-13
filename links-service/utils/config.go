package utils

import (
	"os"
)

type Config struct {
	DBSource       string
	FrontendSource string
}

var (
	ConfigInstance Config
)

func LoadEnvInstance() {
	ConfigInstance = Config{
		DBSource:       os.Getenv("DB_SOURCE"),
		FrontendSource: os.Getenv("FRONTEND_SOURCE"),
	}
}
