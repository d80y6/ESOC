package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/omniguard/ingestion/internal/producer"
	"go.uber.org/zap"
)

type IngestionHTTPServer struct {
	logger   *zap.Logger
	producer *producer.KafkaProducer
}

func NewIngestionHTTPServer(logger *zap.Logger, producer *producer.KafkaProducer) *IngestionHTTPServer {
	return &IngestionHTTPServer{
		logger:   logger,
		producer: producer,
	}
}

type IngestRequest struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (s *IngestionHTTPServer) IngestLog(c echo.Context) error {
	var req IngestRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	eventID := uuid.New().String()

	rawLog := RawLog{
		EventID: eventID,
		Source:  req.Source,
		Type:    req.Type,
		Payload: []byte(req.Payload),
	}

	if err := s.producer.Publish(context.Background(), eventID, rawLog); err != nil {
		s.logger.Error("Failed to publish log to Kafka", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to ingest"})
	}

	s.logger.Info("Ingested HTTP log", zap.String("event_id", eventID))

	return c.JSON(http.StatusOK, map[string]string{
		"event_id": eventID,
		"success":  "true",
	})
}

func (s *IngestionHTTPServer) RegisterRoutes(e *echo.Echo) {
	e.POST("/ingest", s.IngestLog)
}
