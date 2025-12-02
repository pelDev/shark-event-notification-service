package providers

import (
	"fmt"
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
)

type EmailProvider struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	emailFrom    string
}

func NewEmailProvider(host string, port int, username, password, from string) *EmailProvider {
	return &EmailProvider{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		emailFrom:    from,
	}
}

func (p *EmailProvider) Name() string {
	return "Custom Smtp"
}

func (p *EmailProvider) Send(n *domain.Notification) (string, error) {
	if n.Recipient.Email == nil || *n.Recipient.Email == "" {
		return "", fmt.Errorf("email missing for SMS")
	}

	// Simulate email sending
	time.Sleep(100 * time.Millisecond)

	// In production, integrate with email service (SendGrid, AWS SES, etc.)
	email := *n.Recipient.Email
	subject := n.Content.Title
	// body := notification.Content.Body

	fmt.Printf("[SMS] Sending to %s: %s\n", email, subject)

	// return fmt.Sprintf("Email sent to %s with subject: %s", email, subject), nil
	return "", fmt.Errorf("Not Implemented")
}

func (p *EmailProvider) Supports(notificationType domain.NotificationType) bool {
	return notificationType == domain.EmailNotification
}
