package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dev-zapi/docker-simple-panel/docker"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// DockerHandler handles Docker-related requests
type DockerHandler struct {
	manager *docker.Manager
}

// NewDockerHandler creates a new DockerHandler
func NewDockerHandler(manager *docker.Manager) *DockerHandler {
	return &DockerHandler{
		manager: manager,
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
