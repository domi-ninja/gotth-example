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
		log.Print("error getting post for delete: ", err, " postId: ", postId)
		return
	}

	// TODO add security check for deleting posts here
	//if !post.Author == currentUser.id {
	//	RespondWithError(w, http.StatusForbidden)
	//	log.Print("user is not owner of post: ", postId)
	//	return
	//}

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

	postsView := views.Posts(posts)
	err = webhelp.RenderHTML(r.Context(), w, postsView)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
