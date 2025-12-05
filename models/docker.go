package models

// ContainerInfo represents Docker container information
type ContainerInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Image          string `json:"image"`
	State          string `json:"state"`
	Status         string `json:"status"`
	Health         string `json:"health"`
	Created        int64  `json:"created"`
	IsSelf         bool   `json:"is_self"`          // Whether this container is running this application
	ComposeProject string `json:"compose_project"`  // Docker Compose project name
	ComposeService string `json:"compose_service"`  // Docker Compose service name
}

// ContainerOperation represents an operation to perform on a container
type ContainerOperation struct {
	ContainerID string `json:"container_id"`
}

// VolumeInfo represents Docker volume information
type VolumeInfo struct {
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Mountpoint string   `json:"mountpoint"`
	CreatedAt  string   `json:"created_at"`
	Scope      string   `json:"scope"`
	Containers []string `json:"containers"` // List of container IDs using this volume
}

// Response represents a generic API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
