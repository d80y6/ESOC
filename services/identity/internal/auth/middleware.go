package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omniguard/libs/auth"
)

func HasPermission(ctx context.Context, permission string) bool {
	claims, ok := ctx.Value(auth.ClaimsKey).(jwt.MapClaims)
	if !ok {
		return false
	}

	perms, ok := claims["permissions"].([]interface{})
	if !ok {
		return false
	}

	for _, p := range perms {
		if ps, ok := p.(string); ok && ps == permission {
			return true
		}
	}

	return false
}
