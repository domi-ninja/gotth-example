package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
	"github.com/go-chi/chi/v5"
)

func (app *App) HandlePost_PostId_DELETE(w http.ResponseWriter, r *http.Request) {

	postId := chi.URLParam(r, "postId")

	post, err := app.db.GetPostById(r.Context(), postId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound)
		log.Print("hacking, delete non-existing post: ", err, " postId: ", postId)
		return
	}

	user := app.GetCurrentUser(r)
	if user == nil {
		log.Print("hacking, delete while logged out", user, postId)
		RespondWithError(w, http.StatusForbidden)
		return
	}

	if post.UserID != user.ID {
		log.Print("hacking, user is not owner of post ", user, postId)
		RespondWithError(w, http.StatusForbidden)
		return
	}

	app.db.DeletePost(r.Context(), post.ID)

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
