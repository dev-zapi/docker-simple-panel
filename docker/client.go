package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"

	"github.com/dev-zapi/docker-simple-panel/models"
)

const (
	// shortIDLength is the length of the short container ID (12 hex characters)
	shortIDLength = 12
)

// Client wraps the Docker client
type Client struct {
	cli *client.Client
}

// NewClient creates a new Docker client using the Unix socket
func NewClient(socketPath string) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost("unix://"+socketPath),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{cli: cli}, nil
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context) ([]models.ContainerInfo, error) {
	containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []models.ContainerInfo
	for _, container := range containers {
		// Get health status
		health := "none"
		if container.State == "running" {
			inspect, err := c.cli.ContainerInspect(ctx, container.ID)
			if err == nil && inspect.State.Health != nil {
				health = inspect.State.Health.Status
			}
		}

		name := "unknown"
		if len(container.Names) > 0 {
			name = container.Names[0]
			// Remove leading slash from container name
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
		}

		// Extract Docker Compose labels
		composeProject := ""
		composeService := ""
		if container.Labels != nil {
			if project, ok := container.Labels["com.docker.compose.project"]; ok {
				composeProject = project
			}
			if service, ok := container.Labels["com.docker.compose.service"]; ok {
				composeService = service
			}
		}

		result = append(result, models.ContainerInfo{
			ID:             container.ID[:shortIDLength],
			Name:           name,
			Image:          container.Image,
			State:          container.State,
			Status:         container.Status,
			Health:         health,
			Created:        container.Created,
			ComposeProject: composeProject,
			ComposeService: composeService,
		})
	}

	return result, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	return c.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	timeout := 10
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}
	return c.cli.ContainerStop(ctx, containerID, stopOptions)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, containerID string) error {
	timeout := 10
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}
	return c.cli.ContainerRestart(ctx, containerID, stopOptions)
}

// Ping checks if the Docker daemon is accessible
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

// Close closes the Docker client connection
func (c *Client) Close() error {
	return c.cli.Close()
}

// GetContainerInfo gets detailed information about a specific container
func (c *Client) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	inspect, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	health := "none"
	if inspect.State.Health != nil {
		health = inspect.State.Health.Status
	}

	name := inspect.Name
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	created := inspect.Created
	createdTime, _ := time.Parse(time.RFC3339Nano, created)

	// Extract Docker Compose labels
	composeProject := ""
	composeService := ""
	if inspect.Config.Labels != nil {
		if project, ok := inspect.Config.Labels["com.docker.compose.project"]; ok {
			composeProject = project
		}
		if service, ok := inspect.Config.Labels["com.docker.compose.service"]; ok {
			composeService = service
		}
	}

	// Extract restart policy
	var restartPolicy *models.RestartPolicy
	if inspect.HostConfig != nil && inspect.HostConfig.RestartPolicy.Name != "" {
		restartPolicy = &models.RestartPolicy{
			Name:              string(inspect.HostConfig.RestartPolicy.Name),
			MaximumRetryCount: inspect.HostConfig.RestartPolicy.MaximumRetryCount,
		}
	}

	// Extract network information
	networks := make(map[string]models.NetworkInfo)
	if inspect.NetworkSettings != nil && inspect.NetworkSettings.Networks != nil {
		for netName, netConfig := range inspect.NetworkSettings.Networks {
			networks[netName] = models.NetworkInfo{
				NetworkID:  netConfig.NetworkID,
				Gateway:    netConfig.Gateway,
				IPAddress:  netConfig.IPAddress,
				MacAddress: netConfig.MacAddress,
			}
		}
	}

	// Extract port bindings
	var ports []models.PortBinding
	if inspect.NetworkSettings != nil && inspect.NetworkSettings.Ports != nil {
		for port, bindings := range inspect.NetworkSettings.Ports {
			if bindings == nil || len(bindings) == 0 {
				// Port exposed but not bound
				ports = append(ports, models.PortBinding{
					ContainerPort: string(port),
				})
			} else {
				for _, binding := range bindings {
					ports = append(ports, models.PortBinding{
						ContainerPort: string(port),
						HostIP:        binding.HostIP,
						HostPort:      binding.HostPort,
					})
				}
			}
		}
	}

	// Extract mount information
	var mounts []models.MountInfo
	if inspect.Mounts != nil {
		for _, mount := range inspect.Mounts {
			mounts = append(mounts, models.MountInfo{
				Type:        string(mount.Type),
				Source:      mount.Source,
				Destination: mount.Destination,
				Mode:        mount.Mode,
				RW:          mount.RW,
			})
		}
	}

	return &models.ContainerInfo{
		ID:             inspect.ID[:shortIDLength],
		Name:           name,
		Image:          inspect.Config.Image,
		State:          inspect.State.Status,
		Status:         inspect.State.Status,
		Health:         health,
		Created:        createdTime.Unix(),
		ComposeProject: composeProject,
		ComposeService: composeService,
		RestartPolicy:  restartPolicy,
		Env:            inspect.Config.Env,
		Networks:       networks,
		Ports:          ports,
		Mounts:         mounts,
		Hostname:       inspect.Config.Hostname,
	}, nil
}

// ListVolumes lists all Docker volumes with associated container information
func (c *Client) ListVolumes(ctx context.Context) ([]models.VolumeInfo, error) {
	volumes, err := c.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get all containers to build volume-to-container mapping
	containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Build a map of volume name to container IDs
	volumeToContainers := make(map[string][]string)
	for _, container := range containers {
		inspect, err := c.cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			log.Printf("Warning: failed to inspect container %s for volume mapping: %v", container.ID[:shortIDLength], err)
			continue
		}

		// Check mounts for volumes
		for _, mount := range inspect.Mounts {
			if mount.Type == "volume" {
				volumeToContainers[mount.Name] = append(volumeToContainers[mount.Name], container.ID[:shortIDLength])
			}
		}
	}

	var result []models.VolumeInfo
	for _, volume := range volumes.Volumes {
		containers := volumeToContainers[volume.Name]
		if containers == nil {
			containers = []string{}
		}

		result = append(result, models.VolumeInfo{
			Name:       volume.Name,
			Driver:     volume.Driver,
			Mountpoint: volume.Mountpoint,
			CreatedAt:  volume.CreatedAt,
			Scope:      volume.Scope,
			Containers: containers,
		})
	}

	return result, nil
}

// ContainerLogs tail 100 line container logs with follow option
func (c *Client) ContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Timestamps: true,
		Tail:       100,
	}
	
	return c.cli.ContainerLogs(ctx, containerID, options)
}
