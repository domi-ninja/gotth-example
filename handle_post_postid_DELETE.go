package app

import (
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/views"
	"domi.ninja/example-project/internal/db_generated"
	"github.com/go-chi/chi/v5"
)

func (webapp *WebApp) HandlePost_PostId_DELETE(w http.ResponseWriter, r *http.Request) {

	postId := chi.URLParam(r, "postId")

	post, err := webapp.db.GetPostById(r.Context(), postId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound)
		log.Print("error getting post for delete: ", err, " postId: ", postId)
		return
	}

	// TODO security check
	webapp.db.DeletePost(r.Context(), post.ID)

	posts, err := webapp.db.GetPostsPage(r.Context(), db_generated.GetPostsPageParams{
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
