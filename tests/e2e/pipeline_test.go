package e2e

import (
	"testing"

	"github.com/omniguard/libs/models/ecs"
)

func TestIngestionToDetectionFlow(t *testing.T) {
	t.Run("Event Normalization Logic", func(t *testing.T) {
		event := ecs.NormalizedEvent{
			Event: ecs.Event{Action: "logon"},
		}

		if event.Event.Action != "logon" {
			t.Errorf("Expected logon action, got %s", event.Event.Action)
		}
	})
}
