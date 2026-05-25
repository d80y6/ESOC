package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/omniguard/alerting/config"
	"github.com/omniguard/alerting/internal/consumer"
	"github.com/omniguard/alerting/internal/handler"
	"github.com/omniguard/alerting/internal/store"
	"github.com/omniguard/libs/auth"
	"github.com/labstack/echo/v4"
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

	processor := consumer.NewAlertProcessor(cfg.KafkaBrokers, cfg.AlertsTopic, cfg.ConsumerGroupID, db, logger)

	// Authenticator
	authenticator, err := auth.NewAuthenticator(context.Background(), cfg.OIDCProviderURL, cfg.OIDCAudience)
	if err != nil {
		logger.Fatal("failed to create authenticator", zap.Error(err))
	}

	e := echo.New()
	e.Use(authenticator.EchoAuthMiddleware)

	alertHandler := handler.NewAlertHandler(db)
	alertHandler.RegisterRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	logger.Info("Starting Alerting HTTP service", zap.String("port", cfg.HTTPPort))
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)); err != nil {
			logger.Info("Shutting down Alerting HTTP server")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Starting Alerting service...")
	go processor.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Alerting service...")
	cancel()
}
