package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewAuthenticator(ctx context.Context, issuerURL string) (*Authenticator, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
	}
	verifier := provider.Verifier(oidcConfig)

	return &Authenticator{
		provider: provider,
		verifier: verifier,
	}, nil
}

func (a *Authenticator) VerifyToken(ctx context.Context, rawIDToken string) (*jwt.Token, error) {
	// In a real scenario, we use oidc verifier for ID tokens.
	// For Access Tokens, we might need a different approach or use Keycloak's userinfo endpoint.
	// Here we'll implement a generic JWT parser for claims extraction after OIDC verification.

	_, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(rawIDToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token, nil
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
