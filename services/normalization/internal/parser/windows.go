package parser

import (
	"encoding/json"
	"time"

	"github.com/omniguard/libs/models/ecs"
)

type WindowsEventParser struct{}

type WinlogbeatEvent struct {
	Winlog struct {
		EventID       int64  `json:"event_id"`
		ProviderName  string `json:"provider_name"`
		Channel       string `json:"channel"`
		ComputerName  string `json:"computer_name"`
		RecordID      int64  `json:"record_id"`
		Task          string `json:"task"`
		Opcode        string `json:"opcode"`
		Keywords      []string `json:"keywords"`
		Level         string `json:"level"`
		User          struct {
			Identifier string `json:"identifier"`
			Name       string `json:"name"`
			Domain     string `json:"domain"`
			Type       string `json:"type"`
		} `json:"user"`
		EventData map[string]interface{} `json:"event_data"`
	} `json:"winlog"`
	Message string `json:"message"`
}

func (p *WindowsEventParser) Parse(raw []byte) (*ecs.NormalizedEvent, error) {
	var winEvent WinlogbeatEvent
	if err := json.Unmarshal(raw, &winEvent); err != nil {
		return nil, err
	}

	event := &ecs.NormalizedEvent{
		Base: ecs.Base{
			Timestamp: time.Now(),
			Message:   winEvent.Message,
		},
		Event: ecs.Event{
			Dataset:  "windows.event",
			Kind:     "event",
			Action:   winEvent.Winlog.Task,
			Outcome:  "success",
		},
		Host: &ecs.Host{
			Hostname: winEvent.Winlog.ComputerName,
		},
		User: &ecs.User{
			ID:   winEvent.Winlog.User.Identifier,
			Name: winEvent.Winlog.User.Name,
		},
		Metadata: make(map[string]interface{}),
	}

	switch winEvent.Winlog.EventID {
	case 4624:
		event.Event.Category = []string{"authentication"}
		event.Event.Action = "logon"
	case 4688:
		event.Event.Category = []string{"process"}
		event.Event.Action = "process-creation"
	}

	event.Metadata["winlog"] = winEvent.Winlog

	return event, nil
}
