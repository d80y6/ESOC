package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/omniguard/search/config"
	"github.com/omniguard/search/internal/handler"
	"github.com/opensearch-project/opensearch-go/v2"
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

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: cfg.OpenSearchURLs,
	})
	if err != nil {
		logger.Fatal("failed to create opensearch client", zap.Error(err))
	}

	e := echo.New()

	searchHandler := handler.NewSearchHandler(client, logger)
	searchHandler.RegisterRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	logger.Info("Starting Search service", zap.String("port", cfg.HTTPPort))
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)); err != nil {
			logger.Info("Shutting down Search HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Search service...")
}
