package server

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/omniguard/ingestion/api/proto/v1"
	"github.com/omniguard/ingestion/internal/producer"
	"go.uber.org/zap"
)

type IngestionGRPCServer struct {
	pb.UnimplementedIngestionServiceServer
	logger   *zap.Logger
	producer *producer.KafkaProducer
}

func NewIngestionGRPCServer(logger *zap.Logger, producer *producer.KafkaProducer) *IngestionGRPCServer {
	return &IngestionGRPCServer{
		logger:   logger,
		producer: producer,
	}
}

type RawLog struct {
	EventID string `json:"event_id"`
	Source  string `json:"source"`
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

func (s *IngestionGRPCServer) IngestLog(ctx context.Context, req *pb.IngestLogRequest) (*pb.IngestLogResponse, error) {
	eventID := uuid.New().String()

	rawLog := RawLog{
		EventID: eventID,
		Source:  req.Source,
		Type:    req.Type,
		Payload: req.Payload,
	}

	if err := s.producer.Publish(ctx, eventID, rawLog); err != nil {
		s.logger.Error("Failed to publish log to Kafka", zap.Error(err))
		return &pb.IngestLogResponse{Success: false}, err
	}

	s.logger.Info("Ingested gRPC log", zap.String("event_id", eventID))

	return &pb.IngestLogResponse{
		EventId: eventID,
		Success: true,
	}, nil
}
