package config

import (
	"os"
	"time"
)

// Config holds the application configuration
type Config struct {
	APIEndpoint     string
	SourceName      string
	MongoURI        string
	MongoDatabase   string
	MongoCollection string
	FetchInterval   time.Duration
	ServerPort      string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		APIEndpoint:     getEnv("API_ENDPOINT", "https://jsonplaceholder.typicode.com/posts"),
		SourceName:      getEnv("SOURCE_NAME", "placeholder_api"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:   getEnv("MONGO_DATABASE", "logs"),
		MongoCollection: getEnv("MONGO_COLLECTION", "posts"),
		FetchInterval:   getDurationEnv("FETCH_INTERVAL", 5*time.Minute),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
