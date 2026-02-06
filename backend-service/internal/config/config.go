package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config holds all configuration for the application
type Config struct {
	// Application
	AppName    string
	AppVersion string
	Env        string

	// Server
	Host string
	Port int

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Inference Service
	InferenceServiceURL string

	// CORS
	CORSOrigins []string
}

var cfg *Config

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using environment variables")
	}

	cfg = &Config{
		// Application
		AppName:    getEnv("APP_NAME", "BeanSpect Backend Service"),
		AppVersion: getEnv("APP_VERSION", "1.0.0"),
		Env:        getEnv("ENV", "development"),

		// Server
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnvAsInt("PORT", 8080),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "beanspect"),
		DBPassword: getEnv("DB_PASSWORD", "beanspect_secret"),
		DBName:     getEnv("DB_NAME", "beanspect"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Inference Service
		InferenceServiceURL: getEnv("INFERENCE_SERVICE_URL", "http://localhost:8001"),

		// CORS
		CORSOrigins: getEnvAsSlice("CORS_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
	}

	return cfg
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsSlice gets an environment variable as a comma-separated slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}
