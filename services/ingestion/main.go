package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"context"
	"github.com/labstack/echo/v4"
	"github.com/omniguard/ingestion/config"
	"github.com/omniguard/ingestion/internal/producer"
	"github.com/omniguard/ingestion/internal/server"
	"github.com/omniguard/libs/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/omniguard/ingestion/api/proto/v1"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Kafka Producer
	kp := producer.NewKafkaProducer(cfg.KafkaBrokers, cfg.RawLogsTopic, logger)
	defer kp.Close()

	// Authenticator
	authenticator, err := auth.NewAuthenticator(context.Background(), "http://keycloak:8080/realms/omniguard", "omniguard-backend")
	if err != nil {
		logger.Fatal("failed to create authenticator", zap.Error(err))
	}

	// gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen gRPC", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			newCtx, err := authenticator.GRPCAuthInterceptor(ctx)
			if err != nil {
				return nil, err
			}
			return handler(newCtx, req)
		}),
	)
	pb.RegisterIngestionServiceServer(grpcServer, server.NewIngestionGRPCServer(logger, kp))
	reflection.Register(grpcServer)

	logger.Info("Starting Ingestion gRPC service", zap.String("port", cfg.GRPCPort))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// HTTP Server
	e := echo.New()
	e.Use(authenticator.EchoAuthMiddleware)
	httpSrv := server.NewIngestionHTTPServer(logger, kp)
	httpSrv.RegisterRoutes(e)

	logger.Info("Starting Ingestion HTTP service", zap.String("port", cfg.HTTPPort))
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)); err != nil {
			logger.Info("Shutting down HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Ingestion service...")
	grpcServer.GracefulStop()
}
