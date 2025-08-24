package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// HandleLogin_POST handles user login
func (app *App) HandleLogin_POST(w http.ResponseWriter, r *http.Request) {
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

	// Get user by email
	user, err := app.db.GetUserByEmailWithPassword(r.Context(), email)
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
	if !app.CheckPasswordHash(password, user.PasswordHash) {
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
