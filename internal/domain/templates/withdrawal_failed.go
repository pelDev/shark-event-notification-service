package templates

import "fmt"

type WithdrawalFailedData struct {
	Amount      string `json:"amount"`
	ReferenceID string `json:"reference_id"`
	Destination string `json:"destination"`
	Date        string `json:"date"`
	Reason      string `json:"reason"`
}

func (o *WithdrawalFailedData) isEmailTemplateData() {}

func (o *WithdrawalFailedData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *WithdrawalFailedData) GetPreHeader() *string {
	return nil
}
