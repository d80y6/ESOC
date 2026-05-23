package consumer

import (
	"context"
	"encoding/json"

	"github.com/omniguard/normalization/internal/parser"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Worker struct {
	reader   *kafka.Reader
	writer   *kafka.Writer
	logger   *zap.Logger
	parsers  map[string]parser.Parser
}

func NewWorker(brokers []string, rawTopic, normTopic, groupID string, logger *zap.Logger) *Worker {
	return &Worker{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   rawTopic,
			GroupID: groupID,
		}),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    normTopic,
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
		parsers: map[string]parser.Parser{
			"json":    &parser.JSONParser{},
			"syslog":  &parser.SyslogParser{},
			"windows": &parser.WindowsEventParser{},
		},
	}
}

type RawLog struct {
	EventID string `json:"event_id"`
	Source  string `json:"source"`
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

func (w *Worker) Start(ctx context.Context) {
	for {
		m, err := w.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			w.logger.Error("Failed to read message", zap.Error(err))
			continue
		}

		var raw RawLog
		if err := json.Unmarshal(m.Value, &raw); err != nil {
			w.logger.Error("Failed to unmarshal raw log", zap.Error(err))
			continue
		}

		p, ok := w.parsers[raw.Type]
		if !ok {
			p = w.parsers["json"] // Default
		}

		normalized, err := p.Parse(raw.Payload)
		if err != nil {
			w.logger.Error("Failed to parse log", zap.String("event_id", raw.EventID), zap.Error(err))
			continue
		}

		normalized.Event.ID = raw.EventID

		val, _ := json.Marshal(normalized)
		err = w.writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(raw.EventID),
			Value: val,
		})
		if err != nil {
			w.logger.Error("Failed to write normalized event", zap.Error(err))
		}
	}
}

func (w *Worker) Close() {
	w.reader.Close()
	w.writer.Close()
}
