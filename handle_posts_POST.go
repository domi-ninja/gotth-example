package app

import (
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
	"github.com/google/uuid"
)

func (app *App) HandlePosts_POST(w http.ResponseWriter, r *http.Request) {

	user := app.GetCurrentUser(r)
	if user == nil {
		log.Print("hacking, trying to create post while not logged in", r)
		RespondWithError(w, http.StatusForbidden)
		return
	}

	post := db_generated.Post{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	}

	app.db.CreatePost(r.Context(), db_generated.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Title:     post.Title,
		Body:      post.Body,
		UserID:    post.UserID,
	})

	posts, err := app.db.GetPostsPage(r.Context(), db_generated.GetPostsPageParams{
		PagingOffset: 0,
		PageSize:     10,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		log.Print("error getting posts: ", err)
		return
	}

	postsView := views.PostsView(posts, user)
	err = webhelp.RenderHTML(r.Context(), w, postsView)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
