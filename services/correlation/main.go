package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/omniguard/correlation/config"
	"github.com/omniguard/correlation/internal/compiler"
	"github.com/omniguard/correlation/internal/handler"
	"github.com/omniguard/correlation/internal/matcher"
	"github.com/omniguard/libs/auth"
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

	logger.Info("Starting Correlation engine...")
	go engine.Start(ctx)

	// API for Rule Management
	authenticator, err := auth.NewAuthenticator(context.Background(), "http://keycloak:8080/realms/omniguard", "omniguard-backend")
	if err != nil {
		logger.Fatal("failed to create authenticator", zap.Error(err))
	}

	e := echo.New()
	e.Use(authenticator.EchoAuthMiddleware)

	ruleHandler := handler.NewRuleHandler(logger)
	ruleHandler.RegisterRoutes(e)

	logger.Info("Starting Correlation API service", zap.String("port", "8084"))
	go func() {
		if err := e.Start(":8084"); err != nil {
			logger.Info("Shutting down Correlation API server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Correlation engine...")
	cancel()
}
