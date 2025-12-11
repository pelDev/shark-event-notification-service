package templates

import "fmt"

type AdminNewSignupData struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (o *AdminNewSignupData) isEmailTemplateData() {}

func (o *AdminNewSignupData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *AdminNewSignupData) GetPreHeader() *string {
	return nil
}
