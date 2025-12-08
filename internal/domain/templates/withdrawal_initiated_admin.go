package templates

import "fmt"

type WithdrawalInitiatedAdminData struct {
	Amount      string `json:"amount"`
	ReferenceID string `json:"reference_id"`
	Destination string `json:"destination"`
	Date        string `json:"date"`

	Name  string `json:"name"`
	Email string `json:"email"`
}

func (o *WithdrawalInitiatedAdminData) isEmailTemplateData() {}

func (o *WithdrawalInitiatedAdminData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *WithdrawalInitiatedAdminData) GetPreHeader() *string {
	return nil
}
