package templates

import "fmt"

type WelcomeData struct {
	Name string `json:"name"`
}

func (o *WelcomeData) isEmailTemplateData() {}

func (o *WelcomeData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *WelcomeData) GetPreHeader() *string {
	return nil
}
