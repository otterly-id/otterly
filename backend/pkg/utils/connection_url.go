package utils

import (
	"fmt"
	"os"
)

func ConnectionURLBuilder(n string) (string, error) {
	var url string

	switch n {
	case "neon":
		{
			url = os.Getenv("DB_URL")
			if url == "" {
				return "", fmt.Errorf("DB_URL environment variable is not set")
			}
		}

	case "postgres":
		{
			host := os.Getenv("DB_HOST")
			port := os.Getenv("DB_PORT")
			user := os.Getenv("DB_USER")
			password := os.Getenv("DB_PASSWORD")
			dbname := os.Getenv("DB_NAME")
			sslmode := os.Getenv("DB_SSL_MODE")

			if host == "" || port == "" || user == "" || password == "" || dbname == "" {
				return "", fmt.Errorf("missing required database environment variables")
			}

			url = fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, dbname, sslmode,
			)
		}

	case "server":
		{
			host := os.Getenv("SERVER_HOST")
			port := os.Getenv("SERVER_PORT")

			if host == "" || port == "" {
				return "", fmt.Errorf("missing required server environment variables")
			}

			url = fmt.Sprintf("%s:%s", host, port)
		}

	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	return url, nil
}
