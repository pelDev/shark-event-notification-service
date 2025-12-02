package events

import (
	"encoding/json"
	"fmt"
)

type KafkaEvent struct {
	ID        string          `json:"event_id"`
	AggID     string          `json:"aggregate_id"`
	EventType string          `json:"event_type"`
	Version   int             `json:"version"`
	Payload   json.RawMessage `json:"payload"`
}

func (e *KafkaEvent) Validate() error {
	const supportedVersion = 1 // TODO: Move to config

	if e.Version != supportedVersion {
		return fmt.Errorf("unsupported event version: %d, only version %d is supported", e.Version, supportedVersion)
	}

	return nil
}
