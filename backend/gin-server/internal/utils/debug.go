package utils

// isDebugEnabled checks if debug logging should be enabled based on environment or configuration.
func IsDebugEnabled() bool {
	appEnv := GetEnv("APP_ENV", "production") // Default to production if not set
	return appEnv == "local" || appEnv == "development"
}
