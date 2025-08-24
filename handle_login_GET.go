package app

import (
	"net/http"

	"domi.ninja/example-project/frontend/components"
	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/webhelp"
)

func (app *App) HandleLogin_GET(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	if _, err := app.GetCurrentUser(r); err == nil {
		// User is already logged in, redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Create header component with toggle dark (no user email for login page)
	toggleDark := components.ToggleDark()
	header := layouts.Header(app.cfg.Site, toggleDark, "")

	// Create login view
	loginView := views.Login()

	// Create master layout with header and view
	component := layouts.Master(loginView, header, app.cfg.Site, app.cfg.Site, app.version)

	err := webhelp.RenderHTML(r.Context(), w, component)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
