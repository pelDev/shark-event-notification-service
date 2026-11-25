package templates

import "fmt"

type InstallmentData struct {
	N      string `json:"n"`      // e.g. "#1"
	Status string `json:"status"` // installment.Status
	DueAt  string `json:"due_at"` // formatted datetime
	Amount string `json:"amount"` // formatted currency
}

type UserPlanActivatedData struct {
	TicketType  string            `json:"ticket_type"`
	EventName   string            `json:"event_name"`
	PlanLink    string            `json:"plan_link"`
	Schedule    []InstallmentData `json:"schedule"`
	GraceEndsAt string            `json:"grace_ends_at"` // formatted datetime
	EventStarts string            `json:"event_starts"`  // formatted datetime
}

func (data *UserPlanActivatedData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (t *UserPlanActivatedData) isEmailTemplateData() {}

func (data *UserPlanActivatedData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("Your plan for %s is now active. View your payment schedule and event details.", data.EventName)
	return &preHeader
}
