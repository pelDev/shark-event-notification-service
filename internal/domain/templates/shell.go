package templates

import (
	"fmt"
	"html/template"
)

type ShellData struct {
	Body template.HTML `json:"body"`
}

func (tS *ShellData) isEmailTemplateData() {}

func (tS *ShellData) GetMessage(emailFrom, email, subject, html string) []byte {
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"Precedence: bulk\r\n"+
			"X-Mailer: Eventor Newsletter\r\n"+
			"List-Unsubscribe: <mailto:unsubscribe@eventor.com.ng?subject=unsubscribe>\r\n"+
			"List-Unsubscribe-Post: List-Unsubscribe=One-Click\r\n"+
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
