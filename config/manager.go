package config

import (
	"sync"
)

// Manager handles runtime configuration updates
type Manager struct {
	mu                   sync.RWMutex
	dockerSocket         string
	disableRegistration  bool
	logLevel             LogLevel
	volumeExplorerImage  string
	sessionMaxTimeout    int
	onDockerSocketChange func(string) error
}

// NewManager creates a new configuration manager
func NewManager(dockerSocket string, disableRegistration bool, logLevel LogLevel, volumeExplorerImage string, sessionMaxTimeout int) *Manager {
	return &Manager{
		dockerSocket:        dockerSocket,
		disableRegistration: disableRegistration,
		logLevel:            logLevel,
		volumeExplorerImage: volumeExplorerImage,
		sessionMaxTimeout:   sessionMaxTimeout,
	}
}

// GetDockerSocket returns the current Docker socket path
func (m *Manager) GetDockerSocket() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.dockerSocket
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
	m.dockerSocket = socketPath
	m.mu.Unlock()

	return nil
}

// GetDisableRegistration returns whether registration is disabled
func (m *Manager) GetDisableRegistration() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.disableRegistration
}

// SetDisableRegistration updates the registration disabled flag
func (m *Manager) SetDisableRegistration(disabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.disableRegistration = disabled
}

// GetLogLevel returns the current log level
func (m *Manager) GetLogLevel() LogLevel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.logLevel
}

// SetLogLevel updates the log level
func (m *Manager) SetLogLevel(level LogLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logLevel = level
}

// GetVolumeExplorerImage returns the volume explorer image
func (m *Manager) GetVolumeExplorerImage() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.volumeExplorerImage
}

// SetVolumeExplorerImage updates the volume explorer image
func (m *Manager) SetVolumeExplorerImage(image string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.volumeExplorerImage = image
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
	return m.sessionMaxTimeout
}

// SetSessionMaxTimeout updates the session max timeout
func (m *Manager) SetSessionMaxTimeout(timeout int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessionMaxTimeout = timeout
}

// SystemConfig represents the system configuration
type SystemConfig struct {
	DockerSocket        string `json:"docker_socket"`
	DisableRegistration bool   `json:"disable_registration"`
	LogLevel            string `json:"log_level"`
	VolumeExplorerImage string `json:"volume_explorer_image"`
	SessionMaxTimeout   int    `json:"session_max_timeout"`
}

// GetSystemConfig returns the current system configuration
func (m *Manager) GetSystemConfig() SystemConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return SystemConfig{
		DockerSocket:        m.dockerSocket,
		DisableRegistration: m.disableRegistration,
		LogLevel:            m.logLevel.String(),
		VolumeExplorerImage: m.volumeExplorerImage,
		SessionMaxTimeout:   m.sessionMaxTimeout,
	}
}
