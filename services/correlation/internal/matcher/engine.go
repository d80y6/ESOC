package matcher

import (
	"context"
	"encoding/json"

	"github.com/omniguard/correlation/internal/compiler"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Engine struct {
	reader   *kafka.Reader
	writer   *kafka.Writer // To 'alerts' topic
	logger   *zap.Logger
	rules    []*compiler.Rule
}

func NewEngine(brokers []string, inputTopic, outputTopic, groupID string, logger *zap.Logger) *Engine {
	return &Engine{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   inputTopic,
			GroupID: groupID,
		}),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    outputTopic,
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}
}

func (e *Engine) LoadRules(rules []*compiler.Rule) {
	e.rules = rules
}

func (e *Engine) Start(ctx context.Context) {
	for {
		m, err := e.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			e.logger.Error("Failed to read normalized event", zap.Error(err))
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			continue
		}

		for _, rule := range e.rules {
			if rule.Condition(event) {
				e.logger.Info("Detection triggered!", zap.String("rule", rule.Name))

				alert := map[string]interface{}{
					"rule_id":   rule.ID,
					"rule_name": rule.Name,
					"event":     event,
					"timestamp": m.Time,
				}

				val, _ := json.Marshal(alert)
				e.writer.WriteMessages(ctx, kafka.Message{
					Value: val,
				})
			}
		}
	}
}
