package app

import (
	"encoding/json"
	"log"
	"net/http"
)

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
