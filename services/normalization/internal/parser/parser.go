package parser

import (
	"encoding/json"
	"time"

	"github.com/omniguard/libs/models/ecs"
)

type Parser interface {
	Parse(raw []byte) (*ecs.NormalizedEvent, error)
}

type JSONParser struct{}

func (p *JSONParser) Parse(raw []byte) (*ecs.NormalizedEvent, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	event := &ecs.NormalizedEvent{
		Base: ecs.Base{
			Timestamp: time.Now(),
		},
		Metadata: data,
	}

	if ts, ok := data["@timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			event.Timestamp = t
		}
	}

	if msg, ok := data["message"].(string); ok {
		event.Message = msg
	}

	return event, nil
}

type SyslogParser struct{}

func (p *SyslogParser) Parse(raw []byte) (*ecs.NormalizedEvent, error) {
	line := string(raw)

	event := &ecs.NormalizedEvent{
		Base: ecs.Base{
			Timestamp: time.Now(),
			Message:   line,
		},
		Event: ecs.Event{
			Dataset: "syslog",
			Kind:    "event",
		},
	}

	return event, nil
}
