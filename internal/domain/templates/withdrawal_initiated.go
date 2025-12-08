package templates

import "fmt"

type WithdrawalInitiatedData struct {
	Amount      string `json:"amount"`       // 120000
	ReferenceID string `json:"reference_id"` // WDL-9383-ABX
	Destination string `json:"destination"`  // "Bank Account •••• 2210"
	Date        string `json:"date"`         // "Dec 2, 2025"
	Mode        string `json:"mode"`         // "manual"
}

func (o *WithdrawalInitiatedData) isEmailTemplateData() {}

func (o *WithdrawalInitiatedData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *WithdrawalInitiatedData) GetPreHeader() *string {
	return nil
}
