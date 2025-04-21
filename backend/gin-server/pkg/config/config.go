package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/constants"
)

// Config holds all configuration for the application
type Config struct {
	AppEnv string
	Port   int
	DBType string // "sqlite", "mongodb"

	SQLiteFile string

	FirebaseProjectID            string
	UseFirebaseEmulator          bool
	FirebaseEmulatorHost         string
	FirebaseGoogleAppCredentials string

	MongoDBURI  string
	MongoDBName string

	EnableDBSeeding bool

	JWTSecret string
}

// Load loads configuration from environment variables
func Load(fileName string) (*Config, error) {

	err := godotenv.Load(fileName)
	if err != nil && os.Getenv("APP_ENV") != constants.EnvProduction { // Only log if not in production, in production env vars might be set directly
		log.Println("Error loading .env file, using environment variables if set")
	}

	config := &Config{
		AppEnv: getEnv("APP_ENV", constants.EnvDevelopment),
		Port:   getEnvAsInt("PORT", 8080),
		DBType: getEnv("DB_TYPE", constants.SQLite),

		SQLiteFile: getEnv("SQLITE_FILE", "./gin-server.db"),

		FirebaseProjectID:            getEnv("FIREBASE_PROJECT_ID", ""),
		UseFirebaseEmulator:          getEnvAsInt("USE_FIREBASE_EMULATOR", 0) == 1,
		FirebaseEmulatorHost:         getEnv("FIREBASE_EMULATOR_HOST", "localhost:8081"),
		FirebaseGoogleAppCredentials: getEnv("FIREBASE_GOOGLE_APP_CREDENTIALS", ""),

		MongoDBURI:  getEnv("MONGO_DB_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "api-server"),

		EnableDBSeeding: getEnvAsInt("ENABLE_DB_SEEDING", 0) == 1,

		JWTSecret: getEnv("JWT_SECRET", constants.DefaultJWTSecret),
	}

	// Validate configuration
	if config.JWTSecret == constants.DefaultJWTSecret && config.AppEnv == constants.EnvProduction {
		return nil, fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	return config, nil
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
