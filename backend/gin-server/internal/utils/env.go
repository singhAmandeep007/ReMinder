package utils

import "os"

// Helper function to get environment variables with a default value.
func GetEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
