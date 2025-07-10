package config

import (
	"os"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configurable parameters, loaded from .env or defaults.
type Config struct {
	Port               string // HTTP port to bind the server
	ExternalAPIBaseURL string // Base URL for external product API
	StartPage          int    // Initial page number for external API
	BatchSize          int    // Number of CSV records per processing batch
	MaxConcurrency     int    // Maximum concurrent threads for processing
}

// LoadConfig reads environment variables (optionally via .env) into Config.
func LoadConfig() (*Config, error) {
	// Load .env file if present
	_ = godotenv.Load()

	return &Config{
		Port:               getEnv("BACKEND_PORT", "8080"),
		ExternalAPIBaseURL: getEnv("EXTERNAL_API_BASE_URL", "https://hackathon-produtos-api.onrender.com/api/produtos"),
		StartPage:          getEnvAsInt("EXTERNAL_API_START_PAGE", 1),
		BatchSize:          getEnvAsInt("BATCH_SIZE", 1000),
		MaxConcurrency:     getEnvAsInt("MAX_CONCURRENCY", runtime.NumCPU()), // default to # of CPUs
	}, nil
}

// getEnv returns the environment variable or a default.
func getEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

// getEnvAsInt parses an integer env var or returns a default.
func getEnvAsInt(key string, defaultVal int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
