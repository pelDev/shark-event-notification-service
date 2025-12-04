package templates

import (
	"fmt"
)

type OtpData struct {
	OtpCode  string `json:"otp_code"`
	ValidFor int    `json:"valid_for"` // e.g. "5 minutes"
}

func (o *OtpData) isEmailTemplateData() {}

func (o *OtpData) GetMessage(emailFrom, email, subject, html string) []byte {
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
