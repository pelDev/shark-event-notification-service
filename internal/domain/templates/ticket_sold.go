package templates

import (
	"fmt"
)

type TicketSoldData struct {
	EventTitle string `json:"event_title"`
	TicketType string `json:"ticket_type"`
	EventID    string `json:"event_id"`
}

func (tS *TicketSoldData) isEmailTemplateData() {}

func (tS *TicketSoldData) GetMessage(emailFrom, email, subject, html string) []byte {
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		emailFrom,
		email,
		subject,
		html,
	)

	return []byte(message)
}

func (tS *TicketSoldData) GetPreHeader() *string {
	return nil
}
