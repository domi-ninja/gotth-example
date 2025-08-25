package app

import (
	"net/http"
)

// HandleLogout_GET handles user logout
func (app *App) HandleLogout_GET(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT cookie
	app.ClearJWTCookie(w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
