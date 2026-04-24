package templates

import "fmt"

type GuestInviteEvent struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at"`
}

type GuestInviteData struct {
	Subject        string           `json:"subject"`
	CustomNote     string           `json:"custom_note"`
	Event          GuestInviteEvent `json:"event"`
	InvitationRole string           `json:"invitation_role"`
	Link           string           `json:"link"`
	Admits         int              `json:"admits"`
}

func (e *GuestInviteData) isEmailTemplateData() {}

func (e *GuestInviteData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (e *GuestInviteData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("You're invited to %s! Accept your invitation before it expires.", e.Event.Name)
	return &preHeader
}
