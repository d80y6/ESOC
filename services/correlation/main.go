package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/omniguard/correlation/config"
	"github.com/omniguard/correlation/internal/compiler"
	"github.com/omniguard/correlation/internal/matcher"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	engine := matcher.NewEngine(cfg.KafkaBrokers, cfg.InputTopic, cfg.OutputTopic, cfg.ConsumerGroupID, logger)

	// Sample rule
	testRule := &compiler.SigmaRule{
		Title: "Suspicious Login",
		ID:    "1",
		Detection: map[string]interface{}{
			"selection": map[string]interface{}{
				"event.action": "logon",
			},
		},
	}
	compiled, _ := compiler.CompileSigma(testRule)
	engine.LoadRules([]*compiler.Rule{compiled})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Starting Prometheus metrics server on :9090")
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	logger.Info("Starting Correlation engine...")
	go engine.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Correlation engine...")
	cancel()
}
