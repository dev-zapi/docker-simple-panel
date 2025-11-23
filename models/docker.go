package models

// ContainerInfo represents Docker container information
type ContainerInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`
	Status  string `json:"status"`
	Health  string `json:"health"`
	Created int64  `json:"created"`
}

// ContainerOperation represents an operation to perform on a container
type ContainerOperation struct {
	ContainerID string `json:"container_id"`
}

// Response represents a generic API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
