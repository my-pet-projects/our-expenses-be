package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config provides application configuration.
type Config struct {
	Port      int
	Mongo     *MongoConfig
	LogLevel  string
	LogToFile bool
}

// MongoConfig provides MongoDb configuration.
type MongoConfig struct {
	URI      string
	Database string
}

var loadEnvFileFn = godotenv.Load

// ProvideConfiguration initializes the configuration from .env file.
func ProvideConfiguration() (*Config, error) {
	if envFileError := loadEnvFileFn(); envFileError != nil {
		return nil, envFileError
	}

	config := &Config{
		Mongo: &MongoConfig{
			URI:      getEnv("MONGO_URI", ""),
			Database: getEnv("MONGO_DATABASE", ""),
		},
		Port:      getEnvAsInt("PORT", 3000),
		LogLevel:  getEnv("LOG_LEVEL", "DEBUG"),
		LogToFile: getEnvAsBool("LOG_FILE_ENABLE", false),
	}

	return config, nil
}

// Helper function to read an environment variable or return a default value.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Helper function to read an environment variable into integer or return a default value.
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper function to read an environment variable into a bool or return default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
