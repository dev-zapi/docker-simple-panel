package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dev-zapi/docker-simple-panel/database"
	"github.com/dev-zapi/docker-simple-panel/models"
)

// UserHandler handles user management requests
type UserHandler struct {
	db *database.DB
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// ListUsers handles listing all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list users: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
		Message: "User created successfully",
		Data:    user,
	})
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.db.DeleteUser(id); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, "User not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "User deleted successfully",
	})
}
