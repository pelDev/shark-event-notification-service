package templates

import (
	"encoding/json"
	"fmt"
)

type EmailTemplateData interface {
	isEmailTemplateData()
}

type TicketEmailData struct {
	TicketID   string  `json:"ticket_id"`   // {{ ticket_id }}
	QR         string  `json:"qr"`          // {{ qr }} (base64 image)
	EventID    string  `json:"event_id"`    // {{ event_id }}
	EventTitle string  `json:"event_title"` // {{ event_title }}
	TicketType string  `json:"ticket_type"` // {{ ticket_type }}
	Date       string  `json:"date"`        // {{ date }}
	Amount     float64 `json:"amount"`      // {{ amount }}
}

func (tD *TicketEmailData) isEmailTemplateData() {}

func ParseTemplateData(templateName string, data map[string]interface{}, out *EmailTemplateData) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	switch templateName {
	case "ticket-ready":
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("failed to unmarshal ticket email data: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("unknown template name: %s", templateName)
	}
}
