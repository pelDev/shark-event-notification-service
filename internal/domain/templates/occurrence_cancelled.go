package templates

import "fmt"

type OccurrenceCancelledEvent struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type OccurrenceCancelledOccurrence struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	StartDate     string `json:"start_date"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Duration      string `json:"duration"`
	Capacity      int    `json:"capacity"`
	TicketsSold   int    `json:"tickets_sold"`
	AttendeeCount int    `json:"attendee_count"`
}

type OccurrenceCancelledData struct {
	Event       OccurrenceCancelledEvent      `json:"event"`
	Occurrence  OccurrenceCancelledOccurrence `json:"occurrence"`
	CancelledAt string                        `json:"cancelled_at"`
}

func (e *OccurrenceCancelledData) isEmailTemplateData() {}

func (e *OccurrenceCancelledData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *OccurrenceCancelledData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("Event occurrence '%s' for '%s' has been cancelled. %d ticket holders have been notified.",
		e.Occurrence.Label,
		e.Event.Name,
		e.Occurrence.TicketsSold)
	return &preHeader
}
