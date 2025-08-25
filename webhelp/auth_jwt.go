package webhelp

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token for a user
func (app *Wapp) GenerateJWT(userID, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    app.Cfg.Site.Title,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(app.Cfg.Secrets.JWT_SECRET))
}

// ValidateJWT validates a JWT token and returns the claims
func (app *Wapp) ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Cfg.Secrets.JWT_SECRET), nil
	})

	if err != nil {
		log.Print("hacking, jwt.ParseWithClaims failed", tokenString, err)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// SetJWTCookie sets a JWT token as an HTTP-only cookie
func (app *Wapp) SetJWTCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   !strings.HasPrefix(app.Cfg.Site.AppPath, "http://localhost"), // Use secure only in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// ClearJWTCookie clears the JWT cookie
func (app *Wapp) ClearJWTCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   !strings.HasPrefix(app.Cfg.Site.AppPath, "http://localhost"),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// GetJWTFromCookie extracts JWT token from cookie
func (app *Wapp) GetJWTFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
