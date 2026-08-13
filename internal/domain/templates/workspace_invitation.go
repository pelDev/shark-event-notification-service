package templates

import "fmt"

type WorkspaceInvitationData struct {
	WorkspaceName  string `json:"workspace_name"`
	InviteeEmail   string `json:"invitee_email"`
	InviterName    string `json:"inviter_name"`
	InvitationLink string `json:"invitation_link"`
	ExpiresIn      string `json:"expires_in"`
}

func (o *WorkspaceInvitationData) isEmailTemplateData() {}

func (o *WorkspaceInvitationData) GetMessage(emailFrom, email, subject, html string) []byte {
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

func (o *WorkspaceInvitationData) GetPreHeader() *string {
	return nil
}
