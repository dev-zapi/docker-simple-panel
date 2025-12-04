package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/dev-zapi/docker-simple-panel/models"
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
			ID:             container.ID[:12], // Short ID
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

	return &models.ContainerInfo{
		ID:             inspect.ID[:12],
		Name:           name,
		Image:          inspect.Config.Image,
		State:          inspect.State.Status,
		Status:         inspect.State.Status,
		Health:         health,
		Created:        createdTime.Unix(),
		ComposeProject: composeProject,
		ComposeService: composeService,
	}, nil
}
