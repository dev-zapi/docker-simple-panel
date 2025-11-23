package docker

import (
	"context"
	"log"
	"sync"

	"github.com/dev-zapi/docker-simple-panel/models"
)

// Manager manages Docker client with support for runtime socket path changes
type Manager struct {
	mu         sync.RWMutex
	client     *Client
	socketPath string
}

// NewManager creates a new Docker client manager
func NewManager(socketPath string) (*Manager, error) {
	client, err := NewClient(socketPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		client:     client,
		socketPath: socketPath,
	}, nil
}

// RestartWithSocket restarts the Docker client with a new socket path
func (m *Manager) RestartWithSocket(newSocketPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close existing client
	if m.client != nil {
		if err := m.client.Close(); err != nil {
			log.Printf("Warning: failed to close existing Docker client: %v", err)
		}
	}

	// Create new client with new socket path
	newClient, err := NewClient(newSocketPath)
	if err != nil {
		return err
	}

	// Test connection
	ctx := context.Background()
	if err := newClient.Ping(ctx); err != nil {
		// Close the new client before returning error
		if closeErr := newClient.Close(); closeErr != nil {
			log.Printf("Warning: failed to close new Docker client after failed ping: %v", closeErr)
		}
		return err
	}

	m.client = newClient
	m.socketPath = newSocketPath
	log.Printf("Docker client restarted with socket: %s", newSocketPath)

	return nil
}

// GetSocketPath returns the current socket path
func (m *Manager) GetSocketPath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.socketPath
}

// ListContainers lists all containers
func (m *Manager) ListContainers(ctx context.Context) ([]models.ContainerInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.ListContainers(ctx)
}

// GetContainerInfo gets detailed information about a specific container
func (m *Manager) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.GetContainerInfo(ctx, containerID)
}

// StartContainer starts a container
func (m *Manager) StartContainer(ctx context.Context, containerID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.StartContainer(ctx, containerID)
}

// StopContainer stops a container
func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.StopContainer(ctx, containerID)
}

// RestartContainer restarts a container
func (m *Manager) RestartContainer(ctx context.Context, containerID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.RestartContainer(ctx, containerID)
}

// Ping checks if the Docker daemon is accessible
func (m *Manager) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.Ping(ctx)
}

// Close closes the Docker client connection
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}
