package templates

import (
	"fmt"
)

type ShellData struct {
	Body string `json:"body"`
}

func (tS *ShellData) isEmailTemplateData() {}

func (tS *ShellData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (tS *ShellData) GetPreHeader() *string {
	return nil
}
