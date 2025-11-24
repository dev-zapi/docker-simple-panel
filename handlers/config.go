package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dev-zapi/docker-simple-panel/config"
	"github.com/dev-zapi/docker-simple-panel/database"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// ConfigHandler handles system configuration requests
type ConfigHandler struct {
	configManager *config.Manager
	db            *database.DB
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(configManager *config.Manager, db *database.DB) *ConfigHandler {
	return &ConfigHandler{
		configManager: configManager,
		db:            db,
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
	DisableRegistration *bool   `json:"disable_registration,omitempty"`
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

		// Persist to database
		if err := h.db.SetConfig("docker_socket", *req.DockerSocket); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to persist Docker socket config: "+err.Error())
			return
		}
	}

	// Update registration switch if provided
	if req.DisableRegistration != nil {
		h.configManager.SetDisableRegistration(*req.DisableRegistration)

		// Persist to database
		disableRegStr := strconv.FormatBool(*req.DisableRegistration)
		if err := h.db.SetConfig("disable_registration", disableRegStr); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to persist registration config: "+err.Error())
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
