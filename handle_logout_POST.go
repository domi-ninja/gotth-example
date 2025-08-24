package app

import (
	"encoding/json"
	"net/http"
)

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
