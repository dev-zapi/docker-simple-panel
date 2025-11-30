package config

import (
	"os"
	"strings"
)

// LogLevel represents the logging verbosity level
type LogLevel int

const (
	// LogLevelError only logs errors
	LogLevelError LogLevel = iota
	// LogLevelWarn logs warnings and errors
	LogLevelWarn
	// LogLevelInfo logs basic request information
	LogLevelInfo
	// LogLevelDebug logs detailed request/response information
	LogLevelDebug
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "error"
	case LogLevelWarn:
		return "warn"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	default:
		return "info"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "error":
		return LogLevelError
	case "warn", "warning":
		return LogLevelWarn
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	default:
		return LogLevelInfo
	}
}

// Config holds application configuration
type Config struct {
	ServerPort          string
	DatabasePath        string
	JWTSecret           string
	DockerSocket        string
	DisableRegistration bool
	LogLevel            LogLevel
	StaticPath          string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DatabasePath:        getEnv("DATABASE_PATH", "./docker-panel.db"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		DockerSocket:        getEnv("DOCKER_SOCKET", "/var/run/docker.sock"),
		DisableRegistration: getEnvBool("DISABLE_REGISTRATION", false),
		LogLevel:            ParseLogLevel(getEnv("LOG_LEVEL", "info")),
		StaticPath:          getEnv("STATIC_PATH", ""),
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
		lowerValue := strings.ToLower(value)
		return lowerValue == "true" || lowerValue == "1" || lowerValue == "yes"
	}
	return defaultValue
}
