package events

import (
	"encoding/json"
	"fmt"
)

type DomainEvent struct {
	ID        string          `json:"event_id"`
	AggID     string          `json:"aggregate_id"`
	EventType string          `json:"event_type"`
	Version   int             `json:"version"`
	Payload   json.RawMessage `json:"payload"`
}

func (e *DomainEvent) Validate() error {
	const supportedVersion = 1 // TODO: Move to config

	if e.Version != supportedVersion {
		return fmt.Errorf("unsupported event version: %d, only version %d is supported", e.Version, supportedVersion)
	}

	if e.EventType != "notification.requested" {
		return fmt.Errorf("unsupported event type: %s", e.EventType)
	}

	return nil
}
