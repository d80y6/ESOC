package consumer

import (
	"context"
	"encoding/json"

	"github.com/omniguard/alerting/internal/store"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AlertProcessor struct {
	reader *kafka.Reader
	db     *gorm.DB
	logger *zap.Logger
}

func NewAlertProcessor(brokers []string, topic, groupID string, db *gorm.DB, logger *zap.Logger) *AlertProcessor {
	return &AlertProcessor{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		db:     db,
		logger: logger,
	}
}

type KafkaAlert struct {
	RuleID   string                 `json:"rule_id"`
	RuleName string                 `json:"rule_name"`
	Event    map[string]interface{} `json:"event"`
}

func (p *AlertProcessor) Start(ctx context.Context) {
	for {
		m, err := p.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			p.logger.Error("Failed to read alert", zap.Error(err))
			continue
		}

		var ka KafkaAlert
		if err := json.Unmarshal(m.Value, &ka); err != nil {
			continue
		}

		eventJSON, _ := json.Marshal(ka.Event)
		alert := store.Alert{
			RuleID:    ka.RuleID,
			RuleName:  ka.RuleName,
			Status:    "open",
			EventData: string(eventJSON),
		}

		if err := p.db.Create(&alert).Error; err != nil {
			p.logger.Error("Failed to store alert in DB", zap.Error(err))
		} else {
			p.logger.Info("Alert persisted", zap.Uint("id", alert.ID))
		}
	}
}
