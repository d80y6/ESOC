package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/omniguard/normalization/config"
	"github.com/omniguard/normalization/internal/consumer"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	worker := consumer.NewWorker(cfg.KafkaBrokers, cfg.RawLogsTopic, cfg.NormalizedTopic, cfg.ConsumerGroupID, logger)
	defer worker.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Starting Normalization worker...")
	go worker.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Normalization service...")
	cancel()
}
