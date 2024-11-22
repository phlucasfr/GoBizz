package util

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource  string
	MasterKey string
	RedisPort string
	RedisHost string
}

var (
	config     Config
	configOnce sync.Once
)

func GetConfig() *Config {
	configOnce.Do(func() {
		if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "production" {
			config = Config{
				DBSource:  os.Getenv("DB_SOURCE"),
				MasterKey: os.Getenv("MASTER_KEY"),
				RedisPort: os.Getenv("REDIS_PORT"),
				RedisHost: os.Getenv("REDIS_HOST"),
			}
		} else {
			viper.AddConfigPath(".")
			viper.SetConfigName(".env")
			viper.SetConfigType("env")
			viper.AutomaticEnv()

			err := viper.ReadInConfig()
			if err != nil {
				log.Fatalf("Error loading configuration: %v", err)
			}

			err = viper.Unmarshal(&config)
			if err != nil {
				log.Fatalf("Error unmarshalling configuration: %v", err)
			}
		}

		if config.DBSource == "" || config.MasterKey == "" || config.RedisPort == "" || config.RedisHost == "" {
			log.Fatalf("Missing required configuration values")
		}
	})

	return &config
}
