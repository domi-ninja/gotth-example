package app

import (
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
)

func (app *App) Handle_me_password_POST(w http.ResponseWriter, r *http.Request) {
	user := webhelp.GetClaimsFromContext(r.Context())
	if user.UserID == "" {
		log.Print("User not found in context")
		RespondWithError(w, http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Print("Error parsing form:", err)
		RespondWithError(w, http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	// Validate input
	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		log.Print("All password fields are required")
		http.Error(w, "All password fields are required", http.StatusBadRequest)
		return
	}

	if newPassword != confirmPassword {
		log.Print("New passwords do not match")
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	if len(newPassword) < 8 {
		log.Print("New password must be at least 8 characters")
		http.Error(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Get current user from database to verify current password
	dbUser, err := app.db.GetUserById(r.Context(), user.UserID)
	if err != nil {
		log.Print("Error getting user from database:", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Verify current password
	if !app.CheckPasswordHash(currentPassword, dbUser.PasswordHash) {
		log.Print("Current password is incorrect")
		http.Error(w, "Current password is incorrect", http.StatusBadRequest)
		return
	}

	// Hash new password
	hashedPassword, err := app.HashPassword(newPassword)
	if err != nil {
		log.Print("Error hashing password:", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Update user password in database
	err = app.db.UpdateUserPassword(r.Context(), db_generated.UpdateUserPasswordParams{
		PasswordHash: hashedPassword,
		UpdatedAt:    time.Now(),
		ID:           user.UserID,
	})

	if err != nil {
		log.Print("Error updating user password:", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Redirect back to profile page with success message
	http.Redirect(w, r, "/me?password_success=1", http.StatusSeeOther)
}
