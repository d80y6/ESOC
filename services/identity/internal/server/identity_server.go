package server

import (
	"context"

	pb "github.com/omniguard/identity/api/proto/v1"
	"github.com/omniguard/identity/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
	authenticator *auth.Authenticator
}

func NewIdentityServer(auth *auth.Authenticator) *IdentityServer {
	return &IdentityServer{
		authenticator: auth,
	}
}

func (s *IdentityServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := s.authenticator.VerifyToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	tenantID, _ := claims["tenant_id"].(string)
	userID, _ := claims["sub"].(string)

	return &pb.ValidateTokenResponse{
		Valid:    true,
		TenantId: tenantID,
		UserId:   userID,
	}, nil
}
