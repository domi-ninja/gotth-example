package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/components"
	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/webhelp"
)

func (app *App) Handle_me_VIEW(w http.ResponseWriter, r *http.Request) {

	claims := webhelp.GetClaimsFromContext(r.Context())

	// Check for success messages
	var successMessage string
	var passwordSuccessMessage string

	if r.URL.Query().Get("success") == "1" {
		successMessage = "Profile updated successfully!"
	}
	if r.URL.Query().Get("password_success") == "1" {
		passwordSuccessMessage = "Password changed successfully!"
	}

	// Create header component with toggle dark and user email
	toggleDark := components.ToggleDark()
	header := layouts.Header(app.Cfg.Site, toggleDark, claims.Email)

	user, err := app.db.GetUserById(r.Context(), claims.UserID)
	if err != nil {
		log.Print(err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Create me view
	view := views.MeView(user, successMessage, passwordSuccessMessage)

	// Create master layout with header and view
	component := layouts.Master(view, header, app.Cfg.Site, app.Cfg.Site, app.version)

	err = webhelp.RenderHTML(r.Context(), w, component)
	if err != nil {
		log.Print(err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

}
