package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dev-zapi/docker-simple-panel/docker"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// DockerHandler handles Docker-related requests
type DockerHandler struct {
	client *docker.Client
}

// NewDockerHandler creates a new DockerHandler
func NewDockerHandler(client *docker.Client) *DockerHandler {
	return &DockerHandler{
		client: client,
	}
}

// ListContainers handles listing all containers
func (h *DockerHandler) ListContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := h.client.ListContainers(r.Context())
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

	container, err := h.client.GetContainerInfo(r.Context(), containerID)
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

	if err := h.client.StartContainer(r.Context(), containerID); err != nil {
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

	if err := h.client.StopContainer(r.Context(), containerID); err != nil {
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

	if err := h.client.RestartContainer(r.Context(), containerID); err != nil {
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
	if err := h.client.Ping(r.Context()); err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Docker daemon not accessible")
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Docker daemon is accessible",
	})
}
