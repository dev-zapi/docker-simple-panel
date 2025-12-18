package config

import (
	"sync"
)

// Manager handles runtime configuration updates
type Manager struct {
	mu                   sync.RWMutex
	config               *Config
	onDockerSocketChange func(string) error
}

// NewManager creates a new configuration manager
func NewManager(cfg *Config) *Manager {
	return &Manager{
		config: cfg,
	}
}

// GetDockerSocket returns the current Docker socket path
func (m *Manager) GetDockerSocket() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Docker.Socket
}

// SetDockerSocket updates the Docker socket path and triggers Docker client restart
func (m *Manager) SetDockerSocket(socketPath string) error {
	m.mu.Lock()
	callback := m.onDockerSocketChange
	m.mu.Unlock()

	// Call the callback outside the lock to avoid deadlocks
	if callback != nil {
		if err := callback(socketPath); err != nil {
			return err
		}
	}

	m.mu.Lock()
	m.config.Docker.Socket = socketPath
	m.mu.Unlock()

	// Save to config file
	return m.config.Save()
}

// GetLogLevel returns the current log level
func (m *Manager) GetLogLevel() LogLevel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.GetLogLevel()
}

// SetLogLevel updates the log level
func (m *Manager) SetLogLevel(level LogLevel) error {
	m.mu.Lock()
	m.config.Logging.Level = level.String()
	m.mu.Unlock()

	// Save to config file
	return m.config.Save()
}

// GetVolumeExplorerImage returns the volume explorer image
func (m *Manager) GetVolumeExplorerImage() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Docker.VolumeExplorerImage
}

// SetVolumeExplorerImage updates the volume explorer image
func (m *Manager) SetVolumeExplorerImage(image string) error {
	m.mu.Lock()
	m.config.Docker.VolumeExplorerImage = image
	m.mu.Unlock()

	// Save to config file
	return m.config.Save()
}

// SetDockerSocketChangeCallback sets the callback for Docker socket changes
func (m *Manager) SetDockerSocketChangeCallback(callback func(string) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onDockerSocketChange = callback
}

// GetSessionMaxTimeout returns the session max timeout in hours
func (m *Manager) GetSessionMaxTimeout() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Server.SessionMaxTimeout
}

// SetSessionMaxTimeout updates the session max timeout
func (m *Manager) SetSessionMaxTimeout(timeout int) error {
	m.mu.Lock()
	m.config.Server.SessionMaxTimeout = timeout
	m.mu.Unlock()

	// Save to config file
	return m.config.Save()
}

// GetUsername returns the configured username
func (m *Manager) GetUsername() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Username
}

// ValidateCredentials validates username and password
func (m *Manager) ValidateCredentials(username, password string) error {
	return m.config.ValidateCredentials(username, password)
}

// SystemConfig represents the system configuration
type SystemConfig struct {
	DockerSocket        string `json:"docker_socket"`
	LogLevel            string `json:"log_level"`
	VolumeExplorerImage string `json:"volume_explorer_image"`
	SessionMaxTimeout   int    `json:"session_max_timeout"`
	Username            string `json:"username"`
}

// GetSystemConfig returns the current system configuration
func (m *Manager) GetSystemConfig() SystemConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return SystemConfig{
		DockerSocket:        m.config.Docker.Socket,
		LogLevel:            m.config.Logging.Level,
		VolumeExplorerImage: m.config.Docker.VolumeExplorerImage,
		SessionMaxTimeout:   m.config.Server.SessionMaxTimeout,
		Username:            m.config.Username,
	}
}
