package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	TenantIDKey contextKey = "tenant_id"
	ClaimsKey   contextKey = "claims"
)

func (a *Authenticator) AuthInterceptor(ctx context.Context) (context.Context, error) {
	tokenStr, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "missing token: %v", err)
	}

	token, err := a.VerifyToken(ctx, tokenStr)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid claims")
	}

	tenantID, _ := claims["tenant_id"].(string)
	if tenantID == "" {
		tenantID = "default"
	}

	newCtx := context.WithValue(ctx, TenantIDKey, tenantID)
	newCtx = context.WithValue(newCtx, ClaimsKey, claims)

	return newCtx, nil
}

func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return "default"
}

func HasPermission(ctx context.Context, permission string) bool {
	claims, ok := ctx.Value(ClaimsKey).(jwt.MapClaims)
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
