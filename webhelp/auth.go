package webhelp

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// TODO: use this?
func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			// check if is api request by checking if request is xhr
			if err != nil {
				if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				}
				return
			}

			if token == nil {
				if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				}
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
