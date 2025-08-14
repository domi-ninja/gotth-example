package app

import (
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
	"github.com/google/uuid"
)

func (a *App) HandlePosts_POST(w http.ResponseWriter, r *http.Request) {

	post := db_generated.Post{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	}

	a.db.CreatePost(r.Context(), db_generated.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Title:     post.Title,
		Body:      post.Body,
		Author:    post.Author,
	})

	posts, err := a.db.GetPostsPage(r.Context(), db_generated.GetPostsPageParams{
		PagingOffset: 0,
		PageSize:     10,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		log.Print("error getting posts: ", err)
		return
	}

	postsView := views.Posts(posts)
	err = postsView.Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
