package templates

import "fmt"

type PromoterInviteData struct {
	EventName string `json:"event_name"`
	Link      string `json:"link"`
	// Optional personal note from the organizer
	Message string `json:"message"`
}

func (e *PromoterInviteData) isEmailTemplateData() {}

func (d *PromoterInviteData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *PromoterInviteData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("You're invited to promote %s. Accept to get your promo link.", e.EventName)
	return &preHeader
}
