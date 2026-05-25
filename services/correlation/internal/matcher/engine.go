package matcher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/omniguard/correlation/internal/compiler"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Engine struct {
	reader *kafka.Reader
	writer *kafka.Writer // To 'alerts' topic
	logger *zap.Logger

	// Optimization: Group rules by product/service logsource to avoid checking every rule
	// for every event.
	rulesByProduct map[string][]*compiler.Rule
	globalRules    []*compiler.Rule
}

func NewEngine(brokers []string, inputTopic, outputTopic, groupID string, logger *zap.Logger) *Engine {
	return &Engine{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          inputTopic,
			GroupID:        groupID,
			ReadBackoffMin: 100 * time.Millisecond,
			ReadBackoffMax: 1 * time.Second,
		}),
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        outputTopic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			Async:        true,
		},
		logger:         logger,
		rulesByProduct: make(map[string][]*compiler.Rule),
	}
}

func (e *Engine) LoadRules(rules []*compiler.Rule) {
	e.rulesByProduct = make(map[string][]*compiler.Rule)
	e.globalRules = nil

	for _, rule := range rules {
		if rule.Logsource != nil {
			if product, ok := rule.Logsource["product"].(string); ok && product != "" {
				e.rulesByProduct[product] = append(e.rulesByProduct[product], rule)
				continue
			}
		}
		e.globalRules = append(e.globalRules, rule)
	}
	e.logger.Info("Rules loaded",
		zap.Int("total_rules", len(rules)),
		zap.Int("global_rules", len(e.globalRules)),
		zap.Int("product_groups", len(e.rulesByProduct)),
	)
}

func (e *Engine) Start(ctx context.Context) {
	e.logger.Info("Correlation engine started")
	for {
		m, err := e.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			kafkaReadErrors.Inc()
			e.logger.Error("Failed to read normalized event", zap.Error(err))
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			e.logger.Error("Failed to unmarshal event", zap.Error(err))
			continue
		}

		startTime := time.Now()
		matched := e.match(event, m.Time)
		duration := time.Since(startTime)

		// Metrics
		eventsProcessed.Inc()
		matchLatency.Observe(duration.Seconds())

		if matched > 0 {
			e.logger.Debug("Event processed", zap.Int("matches", matched), zap.Duration("latency", duration))
		}
	}
}

func (e *Engine) match(event map[string]interface{}, eventTime time.Time) int {
	matches := 0

	// 1. Check product-specific rules
	product := getProductFromEvent(event)
	if rules, ok := e.rulesByProduct[product]; ok {
		matches += e.evaluateRules(rules, event, eventTime)
	}

	// 2. Check global rules
	matches += e.evaluateRules(e.globalRules, event, eventTime)

	return matches
}

func (e *Engine) evaluateRules(rules []*compiler.Rule, event map[string]interface{}, eventTime time.Time) int {
	matches := 0
	for _, rule := range rules {
		if rule.Condition(event) {
			matches++
			e.triggerAlert(rule, event, eventTime)
		}
	}
	return matches
}

func (e *Engine) triggerAlert(rule *compiler.Rule, event map[string]interface{}, eventTime time.Time) {
	e.logger.Info("Detection triggered!", zap.String("rule", rule.Name), zap.String("rule_id", rule.ID))
	detectionsTriggered.WithLabelValues(rule.Name, rule.ID).Inc()

	alert := map[string]interface{}{
		"rule_id":   rule.ID,
		"rule_name": rule.Name,
		"event":     event,
		"timestamp": eventTime.Format(time.RFC3339),
		"metadata": map[string]string{
			"engine": "correlation-v2",
		},
	}

	val, _ := json.Marshal(alert)
	err := e.writer.WriteMessages(context.Background(), kafka.Message{
		Value: val,
	})
	if err != nil {
		kafkaWriteErrors.Inc()
		e.logger.Error("Failed to write alert to Kafka", zap.Error(err))
	}
}

func getProductFromEvent(event map[string]interface{}) string {
	// Attempt to extract product from ECS metadata
	if metadata, ok := event["metadata"].(map[string]interface{}); ok {
		if product, ok := metadata["product"].(string); ok {
			return product
		}
	}
	// Fallback to event.dataset or similar
	if ev, ok := event["event"].(map[string]interface{}); ok {
		if dataset, ok := ev["dataset"].(string); ok {
			return dataset
		}
	}
	return ""
}

func (e *Engine) Close() {
	e.reader.Close()
	e.writer.Close()
}
