package templates

import "fmt"

type PromoterSaleData struct {
	CommissionEarned string `json:"commission_earned"`
	TicketsSold      int    `json:"tickets_sold"`
	Reference        string `json:"reference"`
	OccurredOn       string `json:"occurred_on"`
	Attribution      string `json:"attribution"`
}

func (e *PromoterSaleData) isEmailTemplateData() {}

func (e *PromoterSaleData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *PromoterSaleData) GetPreHeader() *string {
	return nil
}
