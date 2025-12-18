package models

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}
