package templates

import "fmt"

type EventPublishedData struct {
	EventTitle string `json:"event_title"`
	EventID    string `json:"event_id"`
	EventURL   string `json:"event_url"`
}

func (e *EventPublishedData) isEmailTemplateData() {}

func (e *EventPublishedData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *EventPublishedData) GetPreHeader() *string {
	preHeader := "Your event is now live and ready for attendees."
	return &preHeader
}
