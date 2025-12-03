package templates

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type EmailTemplateData interface {
	isEmailTemplateData()
	GetMessage(emailFrom, email, subject, html string) []byte
}

type TicketEmailData struct {
	TicketID   string `json:"ticket_id"`   // {{ ticket_id }}
	QR         string `json:"qr"`          // {{ qr }} (base64 image)
	EventID    string `json:"event_id"`    // {{ event_id }}
	EventTitle string `json:"event_title"` // {{ event_title }}
	TicketType string `json:"ticket_type"` // {{ ticket_type }}
	Date       string `json:"date"`        // {{ date }}
	Amount     string `json:"amount"`      // {{ amount }}
}

func (tD *TicketEmailData) isEmailTemplateData() {}

func (tD *TicketEmailData) GetMessage(emailFrom, email, subject, html string) []byte {
	boundary := fmt.Sprintf("mixed-%d", time.Now().UnixNano())
	cid := "qr@local"

	// Strip data URL prefix if present
	qrBase64 := strings.TrimPrefix(tD.QR, "data:image/png;base64,")

	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/related; boundary=\"%s\"\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: image/png\r\n"+
			"Content-ID: <%s>\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"\r\n"+
			"%s\r\n"+
			"\r\n"+
			"--%s--\r\n",
		emailFrom,
		email,
		subject,
		boundary,
		boundary,
		html,
		boundary,
		cid,
		qrBase64,
		boundary,
	)

	return []byte(message)
}

func ParseTemplateData(templateName string, data map[string]interface{}, out *EmailTemplateData) error {
	fmt.Printf("[ParseTemplateData] templateName: %s, data: %v", templateName, data)

	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	switch templateName {
	case "ticket-ready":
		var result TicketEmailData
		if err := json.Unmarshal(raw, &result); err != nil {
			return fmt.Errorf("failed to unmarshal ticket email data: %w", err)
		}
		*out = &result
		return nil

	default:
		return fmt.Errorf("unknown template name: %s", templateName)
	}
}
