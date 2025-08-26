package webhelp

import (
	"context"
)

type contextKey string

const (
	userClaimsKey contextKey = "context_user_claims"
)

// SetClaimsInContext stores user information in request context
func SetClaimsInContext(ctx context.Context, user Claims) context.Context {
	ctx = context.WithValue(ctx, userClaimsKey, user)
	return ctx
}

// GetUserIDFromContext retrieves user ID from request context
func GetClaimsFromContext(ctx context.Context) Claims {
	user, _ := ctx.Value(userClaimsKey).(Claims)
	// if !ok {
	// 	retur nil
	// 	// log.Print("failed to cast user claims from context ")
	// }
	return user
}
