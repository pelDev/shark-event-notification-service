package templates

import "fmt"

type PromotionApplicationAcceptedData struct {
	EventName  string `json:"event_name"`
	Link       string `json:"link"`
	AcceptedAt string `json:"accepted_at"`
}

func (e *PromotionApplicationAcceptedData) isEmailTemplateData() {}

func (d *PromotionApplicationAcceptedData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *PromotionApplicationAcceptedData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("You're now a promoter for %s. Grab your link and start earning.", e.EventName)
	return &preHeader
}
