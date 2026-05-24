package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	TenantIDKey contextKey = "tenant_id"
	ClaimsKey   contextKey = "claims"
)

type Authenticator struct {
	verifier *oidc.IDTokenVerifier
}

func NewAuthenticator(ctx context.Context, issuerURL string, clientID string) (*Authenticator, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	return &Authenticator{
		verifier: verifier,
	}, nil
}

// VerifyToken verifies the raw token and returns the parsed JWT token
func (a *Authenticator) VerifyToken(ctx context.Context, rawToken string) (*jwt.Token, error) {
	_, err := a.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(rawToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// gRPC Interceptor
func (a *Authenticator) GRPCAuthInterceptor(ctx context.Context) (context.Context, error) {
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
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant_id in token")
	}

	newCtx := context.WithValue(ctx, TenantIDKey, tenantID)
	newCtx = context.WithValue(newCtx, ClaimsKey, claims)

	return newCtx, nil
}

// Echo Middleware
func (a *Authenticator) EchoAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(401, "missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return echo.NewHTTPError(401, "invalid authorization header format")
		}

		tokenStr := parts[1]
		token, err := a.VerifyToken(c.Request().Context(), tokenStr)
		if err != nil {
			return echo.NewHTTPError(401, fmt.Sprintf("invalid token: %v", err))
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		tenantID, _ := claims["tenant_id"].(string)
		if tenantID == "" {
			return echo.NewHTTPError(401, "missing tenant_id in token")
		}

		c.Set(string(TenantIDKey), tenantID)
		c.Set(string(ClaimsKey), claims)

		return next(c)
	}
}

func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

func ExtractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

func GetEchoTenantID(c echo.Context) string {
	if tenantID, ok := c.Get(string(TenantIDKey)).(string); ok {
		return tenantID
	}
	return ""
}
