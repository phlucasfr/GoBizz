package util

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource  string `mapstructure:"DB_SOURCE"`
	MasterKey string `mapstructure:"MASTER_KEY"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisHost string `mapstructure:"REDIS_HOST"`
}

var (
	config     Config
	configOnce sync.Once
)

func GetConfig(path string) *Config {
	configOnce.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
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
	})
	return &config
}
