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
	authenticator, err := auth.NewAuthenticator(context.Background(), cfg.KeycloakURL+"/realms/"+cfg.KeycloakRealm)
	if err != nil {
		logger.Warn("Failed to initialize Keycloak authenticator, using mock mode for development", zap.Error(err))
		// In a real prod environment, we might fatal here.
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()

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
