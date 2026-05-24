package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"context"
	"github.com/labstack/echo/v4"
	"github.com/omniguard/case-management/config"
	"github.com/omniguard/case-management/internal/handler"
	"github.com/omniguard/case-management/internal/store"
	"github.com/omniguard/libs/auth"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	if err := store.AutoMigrate(db); err != nil {
		logger.Fatal("failed to migrate database", zap.Error(err))
	}

	// Authenticator
	authenticator, err := auth.NewAuthenticator(context.Background(), "http://keycloak:8080/realms/omniguard", "omniguard-backend")
	if err != nil {
		logger.Fatal("failed to create authenticator", zap.Error(err))
	}

	e := echo.New()
	e.Use(authenticator.EchoAuthMiddleware)

	caseHandler := handler.NewCaseHandler(db)
	caseHandler.RegisterRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	logger.Info("Starting Case Management service", zap.String("port", cfg.HTTPPort))
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)); err != nil {
			logger.Info("Shutting down Case Management HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Case Management service...")
}
