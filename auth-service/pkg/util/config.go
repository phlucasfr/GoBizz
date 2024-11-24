package util

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource                  string `mapstructure:"DB_SOURCE"`
	MasterKey                 string `mapstructure:"MASTER_KEY"`
	RedisPort                 string `mapstructure:"REDIS_PORT"`
	RedisHost                 string `mapstructure:"REDIS_HOST"`
	TwilioUsername            string `mapstructure:"TWILIO_USERNAME"`
	TwilioPassword            string `mapstructure:"TWILIO_PASSWORD"`
	TwilioVerificationService string `mapstructure:"TWILIO_VERIFICATION_SERVICE"`
}

var (
	config     Config
	configOnce sync.Once
)

func GetConfig(path string) *Config {
	configOnce.Do(func() {
		if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "production" || IsTesting() || os.Getenv("CI_ENV") == "true" {
			config = Config{
				DBSource:                  os.Getenv("DB_SOURCE"),
				MasterKey:                 os.Getenv("MASTER_KEY"),
				RedisPort:                 os.Getenv("REDIS_PORT"),
				RedisHost:                 os.Getenv("REDIS_HOST"),
				TwilioUsername:            os.Getenv("TWILIO_USERNAME"),
				TwilioPassword:            os.Getenv("TWILIO_PASSWORD"),
				TwilioVerificationService: os.Getenv("TWILIO_VERIFICATION_SERVICE"),
			}
		} else {
			if IsTesting() {
				path = "../../"
			}

			log.Println("istesting", IsTesting())
			viper.AddConfigPath(path)
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
	})

	return &config
}
