package app

import (
	"database/sql"
	"log"
	"net/http"
)

// HandleLogin_POST handles user login
func (app *App) HandleLogin_POST(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max memory
		RespondWithHtmlError(w, r, http.StatusBadRequest, "Invalid form data")
		return
	}

	// Extract form fields
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Basic validation
	if email == "" || password == "" {
		RespondWithHtmlError(w, r, http.StatusOK, "Email and password are required")
		return
	}

	// Get user by email
	user, err := app.db.GetUserByEmailWithPassword(r.Context(), email)
	if err == sql.ErrNoRows {
		RespondWithHtmlError(w, r, http.StatusOK, "Invalid email or password")
		return
	}
	if err != nil {
		log.Printf("Database error getting user: %v", err)
		RespondWithHtmlError(w, r, http.StatusOK, "Internal server error")
		return
	}

	// Check password
	if !app.CheckPasswordHash(password, user.PasswordHash) {
		RespondWithHtmlError(w, r, http.StatusOK, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := app.GenerateJWT(user.ID.(string), user.Email)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		RespondWithHtmlError(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set JWT cookie
	app.SetJWTCookie(w, token)

	// redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
