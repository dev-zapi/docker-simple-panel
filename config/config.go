package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	ServerPort         string
	DatabasePath       string
	JWTSecret          string
	DockerSocket       string
	DisableRegistration bool
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DatabasePath:        getEnv("DATABASE_PATH", "./docker-panel.db"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		DockerSocket:        getEnv("DOCKER_SOCKET", "/var/run/docker.sock"),
		DisableRegistration: getEnvBool("DISABLE_REGISTRATION", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
