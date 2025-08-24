package app

import (
	"net/http"

	"domi.ninja/example-project/frontend/components"
	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/webhelp"
)

func (app *App) HandleRegister_GET(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	if _, err := app.GetCurrentUser(r); err == nil {
		// User is already logged in, redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Create header component with toggle dark (no user email for register page)
	toggleDark := components.ToggleDark()
	header := layouts.Header(app.Cfg.Site, toggleDark, "")

	// Create register view
	registerView := views.Register()

	// Create master layout with header and view
	component := layouts.Master(registerView, header, app.Cfg.Site, app.Cfg.Site, app.version)

	err := webhelp.RenderHTML(r.Context(), w, component)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
