package templates

import "fmt"

type PromotionApplicationRejectedData struct {
	EventName string `json:"event_name"`
	Link      string `json:"link"`
	// Optional: the organizer's reason, if they gave one
	Reason string `json:"reason"`
}

func (e *PromotionApplicationRejectedData) isEmailTemplateData() {}

func (d *PromotionApplicationRejectedData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *PromotionApplicationRejectedData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("An update on your application to promote %s.", e.EventName)
	return &preHeader
}
