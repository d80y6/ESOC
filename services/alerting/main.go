package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/omniguard/alerting/config"
	"github.com/omniguard/alerting/internal/consumer"
	"github.com/omniguard/alerting/internal/store"
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
