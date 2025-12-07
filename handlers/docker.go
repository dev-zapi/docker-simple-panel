package handlers

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/dev-zapi/docker-simple-panel/config"
	"github.com/dev-zapi/docker-simple-panel/docker"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// DockerHandler handles Docker-related requests
type DockerHandler struct {
	manager       *docker.Manager
	configManager *config.Manager
}

// NewDockerHandler creates a new DockerHandler
func NewDockerHandler(manager *docker.Manager, configManager *config.Manager) *DockerHandler {
	return &DockerHandler{
		manager:       manager,
		configManager: configManager,
	}
}

// ListContainers handles listing all containers
func (h *DockerHandler) ListContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := h.manager.ListContainers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list containers: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Data:    containers,
	})
}

// GetContainer handles getting a specific container
func (h *DockerHandler) GetContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["id"]

	if containerID == "" {
		respondWithError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	container, err := h.manager.GetContainerInfo(r.Context(), containerID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Container not found: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Data:    container,
	})
}

// StartContainer handles starting a container
func (h *DockerHandler) StartContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["id"]

	if containerID == "" {
		respondWithError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	if err := h.manager.StartContainer(r.Context(), containerID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start container: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Container started successfully",
	})
}

// StopContainer handles stopping a container
func (h *DockerHandler) StopContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["id"]

	if containerID == "" {
		respondWithError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	if err := h.manager.StopContainer(r.Context(), containerID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to stop container: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Container stopped successfully",
	})
}

// RestartContainer handles restarting a container
func (h *DockerHandler) RestartContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["id"]

	if containerID == "" {
		respondWithError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	if err := h.manager.RestartContainer(r.Context(), containerID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to restart container: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Container restarted successfully",
	})
}

// HealthCheck handles health check requests
func (h *DockerHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.Ping(r.Context()); err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Docker daemon not accessible")
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Docker daemon is accessible",
	})
}

// ListVolumes handles listing all Docker volumes
func (h *DockerHandler) ListVolumes(w http.ResponseWriter, r *http.Request) {
	volumes, err := h.manager.ListVolumes(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list volumes: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Data:    volumes,
	})
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins since CORS middleware already handles origin validation
		// and WebSocket connections are authenticated via JWT
		return true
	},
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// StreamContainerLogs handles WebSocket connections for streaming container logs
func (h *DockerHandler) StreamContainerLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["id"]

	if containerID == "" {
		respondWithError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Create context with cancel for cleanup
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Get container logs starting from 30 minutes ago with follow enabled
	logReader, err := h.manager.ContainerLogs(ctx, containerID, true)
	if err != nil {
		conn.WriteJSON(map[string]string{
			"error": "Failed to get container logs: " + err.Error(),
		})
		return
	}
	defer logReader.Close()

	// Create channels for coordination
	done := make(chan struct{})
	var doneOnce sync.Once
	
	// Goroutine to handle WebSocket pings to keep connection alive
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Failed to send ping: %v", err)
					cancel()
					return
				}
			}
		}
	}()

	// Goroutine to handle client disconnection
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("WebSocket read error (client disconnect): %v", err)
				cancel()
				doneOnce.Do(func() { close(done) })
				return
			}
		}
	}()

	// Use Docker's stdcopy to properly demultiplex stdout and stderr streams
	// Create pipes to separate stdout and stderr
	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()
	
	// Start demuxing in a goroutine with context awareness
	demuxErr := make(chan error, 1)
	go func() {
		defer stdoutWriter.Close()
		defer stderrWriter.Close()
		
		// Monitor context cancellation
		errChan := make(chan error, 1)
		go func() {
			_, err := stdcopy.StdCopy(stdoutWriter, stderrWriter, logReader)
			errChan <- err
		}()
		
		select {
		case err := <-errChan:
			if err != nil && err != io.EOF {
				log.Printf("Error demuxing logs: %v", err)
			}
			demuxErr <- err
		case <-ctx.Done():
			// Context cancelled, close pipes to stop demuxing
			stdoutWriter.Close()
			stderrWriter.Close()
			demuxErr <- ctx.Err()
		}
	}()

	// Channel to merge stdout and stderr while preserving order
	logLines := make(chan string, 100)
	
	// Read from stdout
	go func() {
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case logLines <- scanner.Text():
			}
		}
	}()
	
	// Read from stderr
	go func() {
		scanner := bufio.NewScanner(stderrReader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case logLines <- scanner.Text():
			}
		}
	}()
	
	// Stream logs to WebSocket
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-demuxErr:
			// Demux completed, drain remaining logs
			close(logLines)
			for line := range logLines {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
					return
				}
			}
			if err != nil && err != io.EOF && err != context.Canceled {
				select {
				case <-ctx.Done():
					// Don't try to write if context is cancelled
				default:
					conn.WriteJSON(map[string]string{
						"error": "Error reading logs: " + err.Error(),
					})
				}
			}
			return
		case line, ok := <-logLines:
			if !ok {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket write error: %v", err)
				}
				return
			}
		}
	}
}

// ExploreVolumeFiles handles listing files in a volume
func (h *DockerHandler) ExploreVolumeFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volumeName := vars["name"]
	
	if volumeName == "" {
		respondWithError(w, http.StatusBadRequest, "Volume name is required")
		return
	}
	
	// Get path from query parameter, default to root
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}
	
	// Get the volume explorer image from config
	explorerImage := h.configManager.GetVolumeExplorerImage()
	
	files, err := h.manager.ExploreVolumeFiles(r.Context(), volumeName, path, explorerImage)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to explore volume: "+err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Data:    files,
	})
}

// ReadVolumeFile handles reading a file from a volume
func (h *DockerHandler) ReadVolumeFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volumeName := vars["name"]
	
	if volumeName == "" {
		respondWithError(w, http.StatusBadRequest, "Volume name is required")
		return
	}
	
	// Get file path from query parameter
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		respondWithError(w, http.StatusBadRequest, "File path is required")
		return
	}
	
	// Get the volume explorer image from config
	explorerImage := h.configManager.GetVolumeExplorerImage()
	
	content, err := h.manager.ReadVolumeFile(r.Context(), volumeName, filePath, explorerImage)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read file: "+err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Data:    content,
	})
}
