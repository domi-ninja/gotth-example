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

// HandleRegister_POST handles user registration
func (app *App) HandleRegister_POST(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max memory
		app.respondWithError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	// Extract form fields
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Basic validation
	if email == "" || password == "" {
		app.respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if len(password) < 12 {
		app.respondWithError(w, http.StatusBadRequest, "Password must be at least 12 characters")
		return
	}

	// Check if user already exists
	_, err := app.db.GetUserByEmail(r.Context(), email)
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
	hashedPassword, err := app.HashPassword(password)
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
		Email:        email,
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
