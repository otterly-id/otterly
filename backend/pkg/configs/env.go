package configs

import (
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

type AppConfig struct {
	ServerHost  string
	ServerPort  string
	Environment string

	DatabaseURL          string
	MaxDBConnections     int
	MaxDBIdleConnections int
	MaxDBLifetimeSeconds int
}

func LoadAppConfig() (*AppConfig, error) {
	config := &AppConfig{
		ServerHost:  utils.GetEnv("SERVER_HOST", "localhost"),
		ServerPort:  utils.GetEnv("SERVER_PORT", "8080"),
		Environment: utils.GetEnv("ENV", "development"),

		DatabaseURL:          utils.GetEnv("DB_URL", ""),
		MaxDBConnections:     utils.GetEnvAsInt("DB_MAX_CONNECTIONS", 25),
		MaxDBIdleConnections: utils.GetEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
		MaxDBLifetimeSeconds: utils.GetEnvAsInt("DB_MAX_LIFETIME_CONNECTIONS", 300),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(_ *AppConfig) error {
	requiredEnvs := []string{
		"DB_URL",
		"JWT_SECRET",
	}

	return utils.ValidateRequiredEnvs(requiredEnvs)
}

func (c *AppConfig) GetServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}

func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

func (c *AppConfig) GetDatabaseConfig() map[string]interface{} {
	return map[string]interface{}{
		"url":                  c.DatabaseURL,
		"max_connections":      c.MaxDBConnections,
		"max_idle_connections": c.MaxDBIdleConnections,
		"max_lifetime_seconds": c.MaxDBLifetimeSeconds,
	}
}
