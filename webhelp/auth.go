package webhelp

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Wapp struct {
	Cfg *AppConfig
}

// JWT Claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// HashPassword hashes a plain text password using bcrypt
func (app *Wapp) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+app.Cfg.Secrets.PASSWORD_SALT), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func (app *Wapp) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+app.Cfg.Secrets.PASSWORD_SALT))
	return err == nil
}

// RequireAuth middleware to protect routes
func (app *Wapp) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := app.GetJWTFromCookie(r)
		if err != nil {
			log.Print("hacking, trying to access a route with RequireAuth with no cookie ", r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := app.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Print("claims", claims)

		// Add user info to request context
		ctx := r.Context()
		ctx = SetClaimsInContext(ctx, *claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser extracts current user from JWT token in cookie
func (app *Wapp) GetCurrentUser(r *http.Request) *Claims {
	tokenString, err := app.GetJWTFromCookie(r)
	if err != nil {
		// user not logged in
		return nil
	}

	claims, err := app.ValidateJWT(tokenString)
	log.Print("GetCurrentUser", claims, err)

	if err != nil {
		log.Print("hacking, JWT validation failed for ", tokenString, r)
		return nil
	}

	// TODO we cant ban people like this -- check db for users permissions here
	// userFromDb:=

	return claims
}
