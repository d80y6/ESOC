package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaProducer(brokers []string, topic string, logger *zap.Logger) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, key string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
