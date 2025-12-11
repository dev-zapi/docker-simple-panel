package models

// ContainerInfo represents Docker container information
type ContainerInfo struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Image          string            `json:"image"`
	State          string            `json:"state"`
	Status         string            `json:"status"`
	Health         string            `json:"health"`
	Created        int64             `json:"created"`
	IsSelf         bool              `json:"is_self"`          // Whether this container is running this application
	ComposeProject string            `json:"compose_project"`  // Docker Compose project name
	ComposeService string            `json:"compose_service"`  // Docker Compose service name
	Labels         map[string]string `json:"labels,omitempty"` // All container labels
	RestartPolicy  *RestartPolicy    `json:"restart_policy,omitempty"`
	Env            []string          `json:"env,omitempty"`
	Networks       map[string]NetworkInfo `json:"networks,omitempty"`
	Ports          []PortBinding     `json:"ports,omitempty"`
	Mounts         []MountInfo       `json:"mounts,omitempty"`
	Hostname       string            `json:"hostname,omitempty"`
}

// RestartPolicy represents container restart policy
type RestartPolicy struct {
	Name              string `json:"name"`
	MaximumRetryCount int    `json:"maximum_retry_count,omitempty"`
}

// NetworkInfo represents container network information
type NetworkInfo struct {
	NetworkID   string `json:"network_id"`
	Gateway     string `json:"gateway,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
	MacAddress  string `json:"mac_address,omitempty"`
}

// PortBinding represents port mapping
type PortBinding struct {
	ContainerPort string `json:"container_port"`
	HostIP        string `json:"host_ip,omitempty"`
	HostPort      string `json:"host_port,omitempty"`
}

// MountInfo represents volume/bind mount information
type MountInfo struct {
	Type        string `json:"type"` // bind, volume, tmpfs
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode,omitempty"`
	RW          bool   `json:"rw"`
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

// VolumeFileInfo represents a file or directory in a volume
type VolumeFileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"is_directory"`
	Size        int64  `json:"size"`
	Mode        string `json:"mode"`
	ModTime     string `json:"mod_time"`
}

// VolumeFileContent represents file content from a volume
type VolumeFileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
}
