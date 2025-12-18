package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/dev-zapi/docker-simple-panel/config"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	configManager *config.Manager
	jwtSecret     string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(configManager *config.Manager, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		configManager: configManager,
		jwtSecret:     jwtSecret,
	}
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Validate credentials against config
	if err := h.configManager.ValidateCredentials(req.Username, req.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := h.generateToken(req.Username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Login successful",
		Data: models.LoginResponse{
			Token:    token,
			Username: req.Username,
		},
	})
}

// generateToken generates a JWT token for the user
func (h *AuthHandler) generateToken(username string) (string, error) {
	// Get session timeout from config (in hours)
	sessionTimeout := h.configManager.GetSessionMaxTimeout()
	
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(sessionTimeout)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
