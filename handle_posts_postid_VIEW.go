package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/layouts"
	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
	"github.com/go-chi/chi/v5"
)

func (app *App) HandlePost_PostId_VIEW(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	post, err := app.db.GetPostById(r.Context(), postId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound)
		log.Print("error getting post: ", err)
		return
	}

	postSpecificSeoData := app.Cfg.Site
	postSpecificSeoData.Title = post.Title
	if len(post.Body) > 100 {
		postSpecificSeoData.Description = post.Body[:100] + "..."
	} else {
		postSpecificSeoData.Description = post.Body
	}
	// TODO post images
	// postSpecificSeoData.Image = post.Image
	// postSpecificSeoData.Keywords = post.Keywords

	post2 := db_generated.GetPostsPageRow{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Title:     post.Title,
		Body:      post.Body,
		Email:     post.Email,
		UserID:    post.UserID,
	}

	user := app.GetCurrentUser(r)
	postView := views.PostView(post2, user)
	master := layouts.Master(postView, nil, app.Cfg.Site, postSpecificSeoData, app.version)
	webhelp.RenderHTML(r.Context(), w, master)
}
