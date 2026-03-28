package utils

import "os"

// GetEnv returns the environment variable value or a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvOrDefault is an alias for GetEnv
func GetEnvOrDefault(key, defaultValue string) string {
	return GetEnv(key, defaultValue)
}
