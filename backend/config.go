package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Database configuration
	Dsn string
	// Server configuration
	ApiPort uint

	// Logging configuration
	LogLevel    string
	Environment string

	// Authentication configuration
	JWTSecret      string
	BaseURL        string
	AvatarPath     string
	TokenDuration  int // in minutes
	CookieDuration int // in hours
}

func LoadConfig() (*Config, error) {
	config := &Config{
		// Database configuration
		Dsn: getRequiredEnv("DB_DSN"),

		// Server configuration
		ApiPort: getEnvUintOrDefault("API_PORT", 8080),

		// Logging configuration
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),

		// Authentication configuration
		JWTSecret:      getRequiredEnv("JWT_SECRET"),
		BaseURL:        getEnvOrDefault("BASE_URL", "http://localhost:8080"),
		AvatarPath:     getEnvOrDefault("AVATAR_PATH", "./data/avatars"),
		TokenDuration:  getEnvIntOrDefault("TOKEN_DURATION", 60),  // default 60 minutes
		CookieDuration: getEnvIntOrDefault("COOKIE_DURATION", 24), // default 24 hours
	}
	return config, nil
}

// getRequiredEnv gets an environment variable and panics if it's not set
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloatOrDefault(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatValue)
		}
	}
	return defaultValue
}

func getEnvUintOrDefault(key string, defaultValue uint) uint {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint(intValue)
		}
	}
	return defaultValue
}

func getEnvUintOrPanic(key string) uint64 {
	value := getRequiredEnv(key)
	if uintValue, err := strconv.ParseUint(value, 10, 32); err == nil {
		return uintValue
	}
	panic("Environment variable " + key + " must be a valid unsigned integer")
}
