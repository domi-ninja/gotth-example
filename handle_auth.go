package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/internal/db_generated"
	"github.com/google/uuid"
)

// RegisterRequest represents the JSON body for user registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the JSON body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the JSON response for auth operations
type AuthResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	User    *UserResponse `json:"user,omitempty"`
}

// UserResponse represents user data in responses (no password)
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// HandleRegister_POST handles user registration
func (app *App) HandleRegister_POST(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		app.respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if len(req.Password) < 12 {
		app.respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Check if user already exists
	_, err := app.db.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		app.respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("Database error checking existing user: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Hash the password
	hashedPassword, err := app.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create the user
	userID := uuid.New()
	user, err := app.db.CreateUser(r.Context(), db_generated.CreateUserParams{
		ID:           userID.String(),
		CreatedAt:    time.Now(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate JWT token
	token, err := app.GenerateJWT(user.ID.(string), user.Email)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set JWT cookie
	app.SetJWTCookie(w, token)

	// Respond with success
	response := AuthResponse{
		Success: true,
		Message: "User registered successfully",
		User: &UserResponse{
			ID:        user.ID.(string),
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleLogin_POST handles user login
func (app *App) HandleLogin_POST(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		app.respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user by email
	user, err := app.db.GetUserByEmailWithPassword(r.Context(), req.Email)
	if err == sql.ErrNoRows {
		app.respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	if err != nil {
		log.Printf("Database error getting user: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check password
	if !app.CheckPasswordHash(req.Password, user.PasswordHash) {
		app.respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := app.GenerateJWT(user.ID.(string), user.Email)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set JWT cookie
	app.SetJWTCookie(w, token)

	// Respond with success
	response := AuthResponse{
		Success: true,
		Message: "Login successful",
		User: &UserResponse{
			ID:        user.ID.(string),
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleLogout_POST handles user logout
func (app *App) HandleLogout_POST(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT cookie
	app.ClearJWTCookie(w)

	response := AuthResponse{
		Success: true,
		Message: "Logout successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleMe_GET returns current user information
func (app *App) HandleMe_GET(w http.ResponseWriter, r *http.Request) {
	claims, err := app.GetCurrentUser(r)
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Get full user data from database
	user, err := app.db.GetUserById(r.Context(), claims.UserID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := AuthResponse{
		Success: true,
		Message: "User data retrieved",
		User: &UserResponse{
			ID:        user.ID.(string),
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to respond with JSON error
func (app *App) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := AuthResponse{
		Success: false,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
