package webhelp

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userEmailKey contextKey = "user_email"
)

// SetUserInContext stores user information in request context
func SetUserInContext(ctx context.Context, userID, email string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, userEmailKey, email)
	return ctx
}

// GetUserIDFromContext retrieves user ID from request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// GetUserEmailFromContext retrieves user email from request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}

// GetCurrentUserEmail safely gets the current user's email, returns empty string if not logged in
func (app *Wapp) GetCurrentUserEmail(r *http.Request) string {
	claims, err := app.GetCurrentUser(r)
	if err != nil {
		return ""
	}
	return claims.Email
}
