package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

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

	// Initialize as empty slice to ensure JSON marshals to [] instead of null
	result := []models.ContainerInfo{}
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
	// Initialize as empty slice to ensure JSON marshals to [] instead of null
	ports := []models.PortBinding{}
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
	// Initialize as empty slice to ensure JSON marshals to [] instead of null
	mounts := []models.MountInfo{}
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

	// Initialize as empty slice to ensure JSON marshals to [] instead of null
	result := []models.VolumeInfo{}
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

	// Parse timestamps once for efficient sorting
	type volumeWithTime struct {
		volume models.VolumeInfo
		time   time.Time
		hasErr bool
	}
	
	volumesWithTime := make([]volumeWithTime, len(result))
	for i, vol := range result {
		t, err := time.Parse(time.RFC3339Nano, vol.CreatedAt)
		volumesWithTime[i] = volumeWithTime{
			volume: vol,
			time:   t,
			hasErr: err != nil,
		}
	}

	// Sort volumes by creation time in descending order (newest first)
	sort.Slice(volumesWithTime, func(i, j int) bool {
		// If parsing fails, put the item with invalid time at the end
		if volumesWithTime[i].hasErr && volumesWithTime[j].hasErr {
			return volumesWithTime[i].volume.Name < volumesWithTime[j].volume.Name // fallback to name sorting
		}
		if volumesWithTime[i].hasErr {
			return false
		}
		if volumesWithTime[j].hasErr {
			return true
		}
		
		// Sort in descending order (newest first)
		return volumesWithTime[i].time.After(volumesWithTime[j].time)
	})

	// Extract sorted volumes
	for i, vwt := range volumesWithTime {
		result[i] = vwt.volume
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
		Tail:       "100",
	}
	
	return c.cli.ContainerLogs(ctx, containerID, options)
}

// ExploreVolumeFiles lists files and directories in a volume path using a temporary container
func (c *Client) ExploreVolumeFiles(ctx context.Context, volumeName, path, explorerImage string) ([]models.VolumeFileInfo, error) {
	// Create a temporary container with the volume mounted
	containerName := fmt.Sprintf("volume-explorer-%s-%d", volumeName, time.Now().Unix())
	
	// Create container config
	config := &container.Config{
		Image: explorerImage,
		Cmd:   []string{"ls", "-la", "--full-time", "/volume" + path},
	}
	
	hostConfig := &container.HostConfig{
		Binds: []string{volumeName + ":/volume:ro"}, // Mount as read-only
	}
	
	// Create the container
	resp, err := c.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary container: %w", err)
	}
	
	// Ensure container is removed on exit
	defer func() {
		removeCtx := context.Background()
		c.cli.ContainerRemove(removeCtx, resp.ID, types.ContainerRemoveOptions{Force: true})
	}()
	
	// Start the container
	if err := c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start temporary container: %w", err)
	}
	
	// Wait for container to finish
	statusCh, errCh := c.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("error waiting for container: %w", err)
		}
	case <-statusCh:
	}
	
	// Get container logs (which contains the ls output)
	logReader, err := c.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logReader.Close()
	
	// Use stdcopy to properly demux Docker streams
	var stdout, stderr strings.Builder
	if _, err := stdcopy.StdCopy(&stdout, &stderr, logReader); err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read container output: %w", err)
	}
	
	// Check for errors in stderr
	if stderr.Len() > 0 {
		return nil, fmt.Errorf("command failed: %s", stderr.String())
	}
	
	// Read and parse the output
	// Initialize as empty slice to ensure JSON marshals to [] instead of null
	files := []models.VolumeFileInfo{}
	scanner := bufio.NewScanner(strings.NewReader(stdout.String()))
	
	// Skip the first line (total)
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		
		if firstLine {
			firstLine = false
			if strings.HasPrefix(line, "total") {
				continue
			}
		}
		
		// Parse ls -la output
		// Expected format: mode links owner group size date time timezone filename
		// Example: -rw-r--r-- 1 root root 12 2025-12-07 14:51:49.123456789 +0000 test.txt
		const minFieldCount = 9
		fields := strings.Fields(line)
		if len(fields) < minFieldCount {
			continue
		}
		
		mode := fields[0]
		sizeStr := fields[4]
		dateTime := strings.Join(fields[5:8], " ")
		name := strings.Join(fields[8:], " ")
		
		// Skip . and ..
		if name == "." || name == ".." {
			continue
		}
		
		isDir := strings.HasPrefix(mode, "d")
		size := int64(0)
		if !isDir {
			if s, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
				size = s
			}
		}
		
		// Build full path
		fullPath := path
		if fullPath != "/" && !strings.HasSuffix(fullPath, "/") {
			fullPath += "/"
		}
		if fullPath == "/" {
			fullPath = "/" + name
		} else {
			fullPath += name
		}
		
		files = append(files, models.VolumeFileInfo{
			Name:        name,
			Path:        fullPath,
			IsDirectory: isDir,
			Size:        size,
			Mode:        mode,
			ModTime:     dateTime,
		})
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading container output: %w", err)
	}
	
	return files, nil
}

// ReadVolumeFile reads the content of a file in a volume using a temporary container
func (c *Client) ReadVolumeFile(ctx context.Context, volumeName, filePath, explorerImage string) (*models.VolumeFileContent, error) {
	// Create a temporary container with the volume mounted
	containerName := fmt.Sprintf("volume-reader-%s-%d", volumeName, time.Now().Unix())
	
	// Create container config to read the file
	config := &container.Config{
		Image: explorerImage,
		Cmd:   []string{"cat", "/volume" + filePath},
	}
	
	hostConfig := &container.HostConfig{
		Binds: []string{volumeName + ":/volume:ro"}, // Mount as read-only
	}
	
	// Create the container
	resp, err := c.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary container: %w", err)
	}
	
	// Ensure container is removed on exit with timeout
	defer func() {
		removeCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		c.cli.ContainerRemove(removeCtx, resp.ID, types.ContainerRemoveOptions{Force: true})
	}()
	
	// Start the container
	if err := c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start temporary container: %w", err)
	}
	
	// Wait for container to finish
	statusCh, errCh := c.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("error waiting for container: %w", err)
		}
	case <-statusCh:
	}
	
	// Get container logs (which contains the file content)
	logReader, err := c.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logReader.Close()
	
	// Use stdcopy to properly demux Docker streams
	var stdout, stderr strings.Builder
	if _, err := stdcopy.StdCopy(&stdout, &stderr, logReader); err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}
	
	// Check for errors in stderr
	if stderr.Len() > 0 {
		return nil, fmt.Errorf("failed to read file: %s", stderr.String())
	}
	
	contentStr := stdout.String()
	
	return &models.VolumeFileContent{
		Path:    filePath,
		Content: contentStr,
		Size:    int64(len(contentStr)),
	}, nil
}

func (c *Client) RemoveVolume(ctx context.Context, volumeName string) error {
	return c.cli.VolumeRemove(ctx, volumeName, false)
}
