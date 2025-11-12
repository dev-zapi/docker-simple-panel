package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/dev-zapi/docker-simple-panel/database"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db        *database.DB
	jwtSecret string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *database.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" || req.Nickname == "" {
		respondWithError(w, http.StatusBadRequest, "Username, password, and nickname are required")
		return
	}

	// Create user
	user, err := h.db.CreateUser(req.Username, req.Password, req.Nickname)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, models.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
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

	// Validate credentials
	user, err := h.db.ValidateUser(req.Username, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.Username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Login successful",
		Data: models.LoginResponse{
			Token:    token,
			Username: user.Username,
			Nickname: user.Nickname,
		},
	})
}

// generateToken generates a JWT token for the user
func (h *AuthHandler) generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
