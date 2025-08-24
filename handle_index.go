package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/components"
	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
)

func (app *App) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// Get current user email for header
	userEmail := app.GetCurrentUserEmail(r)

	// Create header component with toggle dark
	toggleDark := components.ToggleDark()
	header := layouts.Header(app.Cfg.Site, toggleDark, userEmail)

	posts, err := app.db.GetPostsPage(r.Context(), db_generated.GetPostsPageParams{
		PagingOffset: 0,
		PageSize:     10,
	})

	if err != nil {
		log.Print("error getting posts: ", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	postsView := views.Posts(posts)

	// Create master layout with header and view
	component := layouts.Master(postsView, header, app.Cfg.Site, app.Cfg.Site, app.version)

	err = webhelp.RenderHTML(r.Context(), w, component)
	if err != nil {
		log.Print("error rendering posts: ", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
