package templates

import "fmt"

type ReminderVenue struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ReminderEvent struct {
	Name            string        `json:"name"`
	Slug            string        `json:"slug"`
	Status          string        `json:"status"`
	Type            string        `json:"type"`
	AccessType      string        `json:"access_type"`
	RequiresPayment bool          `json:"requires_payment"`
	MinTicketPrice  string        `json:"min_ticket_price"`
	Venue           ReminderVenue `json:"venue"`
}

type ReminderOccurrence struct {
	Label     *string `json:"label,omitempty"`
	StartDate string  `json:"start_date"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Duration  string  `json:"duration"`
	Capacity  *int    `json:"capacity,omitempty"`
}

type ReminderEmailData struct {
	Subject    string             `json:"subject"`
	CustomNote string             `json:"custom_note"`
	Event      ReminderEvent      `json:"event"`
	Occurrence ReminderOccurrence `json:"occurrence"`
}

func (e *ReminderEmailData) isEmailTemplateData() {}

func (e *ReminderEmailData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *ReminderEmailData) GetPreHeader() *string {
	preHeader := "Reminder: Your event is coming up soon. Don't miss it!"
	return &preHeader
}
