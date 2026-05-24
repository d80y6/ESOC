package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/omniguard/identity/config"
	"github.com/omniguard/identity/internal/auth"
	"github.com/omniguard/identity/internal/server"
	libsauth "github.com/omniguard/libs/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/omniguard/identity/api/proto/v1"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize Authenticator
	libsAuthenticator, err := libsauth.NewAuthenticator(context.Background(), cfg.KeycloakURL+"/realms/"+cfg.KeycloakRealm, "omniguard-backend")
	if err != nil {
		logger.Fatal("failed to initialize Keycloak authenticator", zap.Error(err))
	}

	authenticator := &auth.Authenticator{Authenticator: libsAuthenticator}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			newCtx, err := authenticator.GRPCAuthInterceptor(ctx)
			if err != nil {
				return nil, err
			}
			return handler(newCtx, req)
		}),
	)

	// Register Identity Server
	identitySrv := server.NewIdentityServer(authenticator)
	pb.RegisterIdentityServiceServer(s, identitySrv)

	reflection.Register(s)

	logger.Info("Starting Identity service", zap.String("port", cfg.GRPCPort))

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Identity service...")
	s.GracefulStop()
}
