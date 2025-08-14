package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"github.com/go-chi/chi/v5"
)

func (a *App) HandlePost_PostId_GET(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	post, err := a.db.GetPostById(r.Context(), postId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound)
		log.Print("error getting post: ", err)
		return
	}

	postView := views.Post(post)
	master := layouts.Master(postView, nil, a.cfg.Site, a.cfg.Site, a.version)
	master.Render(r.Context(), w)
}
