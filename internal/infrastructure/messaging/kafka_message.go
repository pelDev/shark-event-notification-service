package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/commitshark/notification-svc/internal/domain/events"
)

// Infrastructure DTO for Kafka messages

// KafkaNotificationMessage aligns with your NotifyPayload + DomainEvent
type KafkaNotificationMessagePayload struct {
	Type    string `json:"type"` // e.g. "ticket.created"
	Channel string `json:"channel"`
	UserID  string `json:"user_id"` // the user being notified

	Subject  string  `json:"subject"`
	Message  *string `json:"message,omitempty"`  // plain text body
	HTML     *string `json:"html,omitempty"`     // HTML email body
	Template *string `json:"template,omitempty"` // optional template ID

	Data *map[string]any `json:"data,omitempty"` // metadata payload
}

func (m *KafkaNotificationMessagePayload) Validate() error {
	if m.Subject == "" {
		return fmt.Errorf("content.title is required")
	}
	if (m.Message == nil || *m.Message == "") && (m.Data == nil) {
		return fmt.Errorf("content.body and content.data cannot be empty")
	}
	return nil
}

func DecodeNotificationRequestPayload(e *events.KafkaEvent, p *KafkaNotificationMessagePayload) error {
	if e.EventType != "notification.requested" {
		return fmt.Errorf("wrong event type: %s", e.EventType)
	}

	if err := json.Unmarshal(e.Payload, p); err != nil {
		return fmt.Errorf("invalid notification payload: %w", err)
	}

	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}
