package webhelp

import (
	"context"
)

type contextKey string

const (
	userClaimsKey contextKey = "context_user_claims"
)

// SetUserInContext stores user information in request context
func SetUserInContext(ctx context.Context, user Claims) context.Context {
	ctx = context.WithValue(ctx, userClaimsKey, user)
	return ctx
}

// GetUserIDFromContext retrieves user ID from request context
func GetUserFromContext(ctx context.Context) Claims {
	user, _ := ctx.Value(userClaimsKey).(Claims)
	// if !ok {
	// 	retur nil
	// 	// log.Print("failed to cast user claims from context ")
	// }
	return user
}
