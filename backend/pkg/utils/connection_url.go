package utils

import (
	"fmt"
)

func ConnectionURLBuilder(n string) (string, error) {
	var url string

	switch n {
	case "neon":
		{
			url, err := GetEnvRequired("DB_URL")
			if err != nil {
				return "", fmt.Errorf("DB_URL environment variable is not set: %w", err)
			}
			return url, nil
		}

	case "postgres":
		{
			host, err := GetEnvRequired("DB_HOST")
			if err != nil {
				return "", fmt.Errorf("DB_HOST environment variable is not set: %w", err)
			}

			port, err := GetEnvRequired("DB_PORT")
			if err != nil {
				return "", fmt.Errorf("DB_PORT environment variable is not set: %w", err)
			}

			user, err := GetEnvRequired("DB_USER")
			if err != nil {
				return "", fmt.Errorf("DB_USER environment variable is not set: %w", err)
			}

			password, err := GetEnvRequired("DB_PASSWORD")
			if err != nil {
				return "", fmt.Errorf("DB_PASSWORD environment variable is not set: %w", err)
			}

			dbname, err := GetEnvRequired("DB_NAME")
			if err != nil {
				return "", fmt.Errorf("DB_NAME environment variable is not set: %w", err)
			}

			sslmode := GetEnv("DB_SSL_MODE", "require")

			url = fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, dbname, sslmode,
			)

			return url, nil
		}

	case "server":
		{
			host, err := GetEnvRequired("SERVER_HOST")
			if err != nil {
				return "", fmt.Errorf("SERVER_HOST environment variable is not set: %w", err)
			}

			port, err := GetEnvRequired("SERVER_PORT")
			if err != nil {
				return "", fmt.Errorf("SERVER_PORT environment variable is not set: %w", err)
			}

			url = fmt.Sprintf("%s:%s", host, port)
		}

	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	return url, nil
}
