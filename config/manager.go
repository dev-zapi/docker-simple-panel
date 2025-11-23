package config

import (
	"sync"
)

// Manager handles runtime configuration updates
type Manager struct {
	mu                  sync.RWMutex
	dockerSocket        string
	disableRegistration bool
	onDockerSocketChange func(string) error
}

// NewManager creates a new configuration manager
func NewManager(dockerSocket string, disableRegistration bool) *Manager {
	return &Manager{
		dockerSocket:        dockerSocket,
		disableRegistration: disableRegistration,
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
	defer m.mu.Unlock()
	
	// Call the callback to restart Docker client if set
	if m.onDockerSocketChange != nil {
		if err := m.onDockerSocketChange(socketPath); err != nil {
			return err
		}
	}
	
	m.dockerSocket = socketPath
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

// SetDockerSocketChangeCallback sets the callback for Docker socket changes
func (m *Manager) SetDockerSocketChangeCallback(callback func(string) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onDockerSocketChange = callback
}

// SystemConfig represents the system configuration
type SystemConfig struct {
	DockerSocket        string `json:"docker_socket"`
	DisableRegistration bool   `json:"disable_registration"`
}

// GetSystemConfig returns the current system configuration
func (m *Manager) GetSystemConfig() SystemConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return SystemConfig{
		DockerSocket:        m.dockerSocket,
		DisableRegistration: m.disableRegistration,
	}
}
