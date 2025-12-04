package templates

type OrganizerEventReminderData struct {
	EventTitle string
	EventDate  string
	TimeLeft   string // "24 hours", "6 hours", "1 hour"
	EventID    string
}

// TODO: Implement EmailTemplateData interface methods
