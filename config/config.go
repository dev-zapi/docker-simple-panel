package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

// Bcrypt hash prefixes
const (
	bcryptPrefix2a = "$2a$"
	bcryptPrefix2b = "$2b$"
	bcryptPrefix2y = "$2y$"
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

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port              string `yaml:"port"`
	JWTSecret         string `yaml:"jwt_secret"`
	SessionMaxTimeout int    `yaml:"session_max_timeout"`
}

// DockerConfig holds Docker-specific configuration
type DockerConfig struct {
	Socket             string `yaml:"socket"`
	VolumeExplorerImage string `yaml:"volume_explorer_image"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string `yaml:"level"`
}

// Config holds application configuration loaded from YAML
type Config struct {
	Username   string         `yaml:"username"`
	Password   string         `yaml:"password"`
	Server     ServerConfig   `yaml:"server"`
	Docker     DockerConfig   `yaml:"docker"`
	Logging    LoggingConfig  `yaml:"logging"`
	StaticPath string         `yaml:"static_path"`
	
	// Runtime fields (not persisted)
	configPath     string
	mu             sync.RWMutex
	hashedPassword string
}

const defaultConfigPath = "./config.yaml"

// LoadConfig loads configuration from YAML file
func LoadConfig() (*Config, error) {
	configPath := getEnv("CONFIG_PATH", defaultConfigPath)
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		log.Printf("Config file not found at %s, creating default config", configPath)
		cfg := getDefaultConfig()
		cfg.configPath = configPath
		
		// Hash the default password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		cfg.Password = string(hashedPassword)
		cfg.hashedPassword = string(hashedPassword)
		
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}
	
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	cfg.configPath = configPath
	
	// Hash password for validation
	if cfg.Password != "" {
		// Check if password is already hashed
		if !isBcryptHash(cfg.Password) {
			// Password is in plain text, hash it
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			cfg.hashedPassword = string(hashedPassword)
			// Update the config file with hashed password
			cfg.Password = string(hashedPassword)
			if err := cfg.Save(); err != nil {
				log.Printf("Warning: failed to save hashed password to config: %v", err)
			}
		} else {
			// Password is already hashed
			cfg.hashedPassword = cfg.Password
		}
	}
	
	return &cfg, nil
}

// isBcryptHash checks if a string is a bcrypt hash
func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, bcryptPrefix2a) ||
		strings.HasPrefix(s, bcryptPrefix2b) ||
		strings.HasPrefix(s, bcryptPrefix2y)
}

// Save saves the current configuration to YAML file
func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(c.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// ValidateCredentials validates username and password against config
func (c *Config) ValidateCredentials(username, password string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if username != c.Username {
		return fmt.Errorf("invalid credentials")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(c.hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	
	return nil
}

// GetLogLevel returns the parsed log level
func (c *Config) GetLogLevel() LogLevel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return ParseLogLevel(c.Logging.Level)
}

// UpdatePassword updates the password and saves to config
func (c *Config) UpdatePassword(newPassword string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	c.Password = string(hashedPassword)
	c.hashedPassword = string(hashedPassword)
	
	return nil
}

// getDefaultConfig returns a default configuration
func getDefaultConfig() *Config {
	return &Config{
		Username: "admin",
		Password: "changeme",
		Server: ServerConfig{
			Port:              "8080",
			JWTSecret:         "your-secret-key-change-in-production",
			SessionMaxTimeout: 24,
		},
		Docker: DockerConfig{
			Socket:             "/var/run/docker.sock",
			VolumeExplorerImage: "ghcr.io/dev-zapi/docker-simple-panel:latest",
		},
		Logging: LoggingConfig{
			Level: "info",
		},
		StaticPath: "",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
