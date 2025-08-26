package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
)

func (app *App) Handle_me_POST(w http.ResponseWriter, r *http.Request) {
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

	email := r.FormValue("email")
	displayName := r.FormValue("display_name")

	// Validate input
	if email == "" {
		log.Print("Email is required")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Prepare update parameters
	var displayNameNull sql.NullString
	if displayName != "" {
		displayNameNull = sql.NullString{String: displayName, Valid: true}
	} else {
		displayNameNull = sql.NullString{Valid: false}
	}

	// Update user in database
	err := app.db.UpdateUser(r.Context(), db_generated.UpdateUserParams{
		Email:       email,
		DisplayName: displayNameNull,
		UpdatedAt:   time.Now(),
		ID:          user.UserID,
	})

	if err != nil {
		log.Print("Error updating user:", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Redirect back to profile page with success message
	http.Redirect(w, r, "/me?success=1", http.StatusSeeOther)
}
