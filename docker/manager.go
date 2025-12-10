package docker

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/dev-zapi/docker-simple-panel/models"
)

// ErrSelfOperation is returned when attempting to stop/restart the container running this application
var ErrSelfOperation = errors.New("cannot stop or restart the container running this application")

// Manager manages Docker client with support for runtime socket path changes
type Manager struct {
	mu                   sync.RWMutex
	client               *Client
	socketPath           string
	containerEnvironment ContainerEnvironment
}

// NewManager creates a new Docker client manager
func NewManager(socketPath string) (*Manager, error) {
	client, err := NewClient(socketPath)
	if err != nil {
		return nil, err
	}

	// Detect container environment at startup
	env := DetectContainerEnvironment()
	if env.IsInContainer {
		log.Printf("Running inside container (ID: %s)", env.ContainerID)
	} else {
		log.Println("Running outside container environment")
	}

	return &Manager{
		client:               client,
		socketPath:           socketPath,
		containerEnvironment: env,
	}, nil
}

// IsInContainer returns whether the application is running inside a container
func (m *Manager) IsInContainer() bool {
	return m.containerEnvironment.IsInContainer
}

// GetOwnContainerID returns the container ID of this application if running in a container
func (m *Manager) GetOwnContainerID() string {
	return m.containerEnvironment.ContainerID
}

// isSelfContainer checks if the given container ID matches this application's container
func (m *Manager) isSelfContainer(containerID string) bool {
	if !m.containerEnvironment.IsInContainer || m.containerEnvironment.ContainerID == "" {
		return false
	}

	// Normalize both IDs to short form (12 hex chars) for comparison
	selfID := m.containerEnvironment.ContainerID
	targetID := containerID

	// Handle short IDs (12 hex chars) and full IDs (64 hex chars)
	if len(selfID) > 12 {
		selfID = selfID[:12]
	}
	if len(targetID) > 12 {
		targetID = targetID[:12]
	}

	return strings.EqualFold(selfID, targetID)
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

	containers, err := m.client.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	// Mark self-container
	for i := range containers {
		containers[i].IsSelf = m.isSelfContainer(containers[i].ID)
	}

	return containers, nil
}

// GetContainerInfo gets detailed information about a specific container
func (m *Manager) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	info, err := m.client.GetContainerInfo(ctx, containerID)
	if err != nil {
		return nil, err
	}

	// Mark self-container
	info.IsSelf = m.isSelfContainer(info.ID)

	return info, nil
}

// StartContainer starts a container
func (m *Manager) StartContainer(ctx context.Context, containerID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.StartContainer(ctx, containerID)
}

// StopContainer stops a container
func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	// Check if attempting to stop self
	if m.isSelfContainer(containerID) {
		return ErrSelfOperation
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.StopContainer(ctx, containerID)
}

// RestartContainer restarts a container
func (m *Manager) RestartContainer(ctx context.Context, containerID string) error {
	// Check if attempting to restart self
	if m.isSelfContainer(containerID) {
		return ErrSelfOperation
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.RestartContainer(ctx, containerID)
}

// ListVolumes lists all Docker volumes with container associations
func (m *Manager) ListVolumes(ctx context.Context) ([]models.VolumeInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.ListVolumes(ctx)
}

// RemoveVolume removes a Docker volume by name
func (m *Manager) RemoveVolume(ctx context.Context, volumeName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.client.RemoveVolume(ctx, volumeName)
}

// Ping checks if the Docker daemon is accessible
func (m *Manager) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.Ping(ctx)
}

// ContainerLogs gets container logs starting from 30 minutes ago
func (m *Manager) ContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.ContainerLogs(ctx, containerID, follow)
}

// ExploreVolumeFiles lists files and directories in a volume path
func (m *Manager) ExploreVolumeFiles(ctx context.Context, volumeName, path, explorerImage string) ([]models.VolumeFileInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.ExploreVolumeFiles(ctx, volumeName, path, explorerImage)
}

// ReadVolumeFile reads the content of a file in a volume
func (m *Manager) ReadVolumeFile(ctx context.Context, volumeName, filePath, explorerImage string) (*models.VolumeFileContent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client.ReadVolumeFile(ctx, volumeName, filePath, explorerImage)
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
