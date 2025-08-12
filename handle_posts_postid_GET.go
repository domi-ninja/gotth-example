package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"github.com/go-chi/chi/v5"
)

func (webapp *WebApp) HandlePost_PostId_GET(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	post, err := webapp.db.GetPostById(r.Context(), postId)
	if err != nil {
		respondWithError(w, http.StatusNotFound)
		log.Print("error getting post: ", err)
		return
	}

	postView := views.Post(post)
	master := layouts.Master(postView, nil, webapp.siteCfg.Site, webapp.siteCfg.Site, webapp.buildRandomNumber)
	master.Render(r.Context(), w)
}
