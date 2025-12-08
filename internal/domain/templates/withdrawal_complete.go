package templates

import "fmt"

type WithdrawalCompleteData struct {
	Amount      string `json:"amount"`       // 50000.00
	ReferenceID string `json:"reference_id"` // WDL-231222-XF9
	Destination string `json:"destination"`  // "Bank Account •••• 4421"
	Date        string `json:"date"`         // "Dec 2, 2025"
}

func (o *WithdrawalCompleteData) isEmailTemplateData() {}

func (o *WithdrawalCompleteData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *WithdrawalCompleteData) GetPreHeader() *string {
	return nil
}
