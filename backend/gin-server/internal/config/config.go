package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string
	DBType string // "postgres", "sqlite", or "mongodb"

	SqliteFile string // For SQLite

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string

	MongoDBUri  string
	MongoDBName string

	JWTSecret string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil && os.Getenv("APP_ENV") != "cloud" { // Only log if not in cloud, in cloud env vars might be set directly
		log.Println("Error loading .env file, using environment variables if set")
	}

	return Config{
		AppEnv: os.Getenv("APP_ENV"),
		Port:   os.Getenv("PORT"),
		DBType: os.Getenv("DB_TYPE"),

		SqliteFile: os.Getenv("SQLITE_FILE"),

		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDBName:   os.Getenv("POSTGRES_DB_NAME"),

		MongoDBUri:  os.Getenv("MONGO_DB_URI"),
		MongoDBName: os.Getenv("MONGO_DB_NAME"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
