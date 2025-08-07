package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	config.AddConfigPath("./..")

	config.SetDefault("JWT_EXPIRES_IN", 24)
	config.SetDefault("SERVER_URL", "0.0.0.0:8080")
	config.SetDefault("DB_MAX_CONNECTIONS", 100)
	config.SetDefault("DB_MAX_IDLE_CONNECTIONS", 10)
	config.SetDefault("DB_MAX_LIFETIME_CONNECTIONS", 2)

	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error .env file: %w", err))
	}

	return config
}
