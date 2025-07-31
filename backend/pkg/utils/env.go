package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetEnvRequired(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("required environment variable '%s' is not set", key)
}

func GetEnvAsInt(key string, defaultValue int) int {
	value := GetEnv(key, fmt.Sprintf("%d", defaultValue))
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := GetEnv(key, defaultValue.String())
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return defaultValue
}

func GetEnvAsDurationRequired(key string) (time.Duration, error) {
	value, err := GetEnvRequired(key)
	if err != nil {
		return 0, err
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable '%s' must be a valid duration, got: %s", key, value)
	}
	return duration, nil
}

func HasEnv(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

func ValidateRequiredEnvs(requiredKeys []string) error {
	var missingKeys []string

	for _, key := range requiredKeys {
		if !HasEnv(key) {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingKeys, ", "))
	}

	return nil
}
