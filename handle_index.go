package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/components"
	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
)

func (webapp *WebApp) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// Create header component with toggle dark
	toggleDark := components.ToggleDark()
	header := layouts.Header(webapp.siteCfg.Site, toggleDark)

	posts, err := webapp.db.GetPostsPage(r.Context(), db_generated.GetPostsPageParams{
		PagingOffset: 0,
		PageSize:     10,
	})

	if err != nil {
		log.Print("error getting posts: ", err)
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	postsView := views.Posts(posts)

	// Create master layout with header and view
	component := layouts.Master(postsView, header, webapp.siteCfg.Site, webapp.siteCfg.Site, webapp.buildRandomNumber)

	w.Header().Set("Content-Type", "text/html")
	err = component.Render(r.Context(), w)
	if err != nil {
		log.Print("error rendering posts: ", err)
		respondWithError(w, http.StatusInternalServerError)
		return
	}
}
