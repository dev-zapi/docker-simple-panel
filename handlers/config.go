package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dev-zapi/docker-simple-panel/config"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// ConfigHandler handles system configuration requests
type ConfigHandler struct {
	configManager *config.Manager
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(configManager *config.Manager) *ConfigHandler {
	return &ConfigHandler{
		configManager: configManager,
	}
}

// GetConfig retrieves the current system configuration
func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.configManager.GetSystemConfig()

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Data:    cfg,
	})
}

// UpdateConfigRequest represents configuration update request
type UpdateConfigRequest struct {
	DockerSocket        *string `json:"docker_socket,omitempty"`
	LogLevel            *string `json:"log_level,omitempty"`
	VolumeExplorerImage *string `json:"volume_explorer_image,omitempty"`
	SessionMaxTimeout   *int    `json:"session_max_timeout,omitempty"`
}

// UpdateConfig updates system configuration
func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update Docker socket if provided
	if req.DockerSocket != nil {
		if err := h.configManager.SetDockerSocket(*req.DockerSocket); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update Docker socket: "+err.Error())
			return
		}
	}

	// Update log level if provided
	if req.LogLevel != nil {
		logLevel := config.ParseLogLevel(*req.LogLevel)
		if err := h.configManager.SetLogLevel(logLevel); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update log level: "+err.Error())
			return
		}
	}

	// Update volume explorer image if provided
	if req.VolumeExplorerImage != nil {
		if err := h.configManager.SetVolumeExplorerImage(*req.VolumeExplorerImage); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update volume explorer image: "+err.Error())
			return
		}
	}

	// Update session max timeout if provided
	if req.SessionMaxTimeout != nil {
		if err := h.configManager.SetSessionMaxTimeout(*req.SessionMaxTimeout); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update session max timeout: "+err.Error())
			return
		}
	}

	// Return updated configuration
	cfg := h.configManager.GetSystemConfig()
	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Configuration updated successfully",
		Data:    cfg,
	})
}
