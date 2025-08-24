package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the application

var (
	// DBConfig holds the database configuration
	DBConfig = struct {
		Host     string
		User     string
		Password string
		Port     string
		DBName   string
	}{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASS", ""),
		Port:     getEnv("DB_PORT", "3306"),
		DBName:   getEnv("DB_NAME", "sample"),
	}

	// Port holds the server port configuration
	Port   = getEnv("PORT", "8080")
	UseJWT = getEnv("USE_JWT", "true")
)

func getEnv(key, defaultValue string) string {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
		log.Printf("Warning: Error converting %s to int: %v. Using default value.", key, err)
	}
	return defaultValue
}
