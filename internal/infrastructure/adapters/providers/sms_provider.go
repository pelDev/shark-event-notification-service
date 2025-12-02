package providers

import (
	"fmt"

	"github.com/commitshark/notification-svc/internal/domain"
)

type SMSProvider struct {
}

func NewSMSProvider() *SMSProvider {
	return &SMSProvider{}
}

func (p *SMSProvider) Name() string {
	return "Sample Sms Provider"
}

func (p *SMSProvider) Supports(notificationType domain.NotificationType) bool {
	return notificationType == domain.SMSNotification
}

func (p *SMSProvider) Send(n *domain.Notification) (string, error) {
	if n.Recipient.Phone == nil || *n.Recipient.Phone == "" {
		return "", fmt.Errorf("phone number missing for SMS")
	}

	if n.Content.Body == nil || *n.Content.Body == "" {
		return "", fmt.Errorf("message content for SMS")
	}

	// Example; replace with actual API call
	fmt.Printf("[SMS] Sending to %s: %s\n", *n.Recipient.Phone, *n.Content.Body)

	// Return a fake provider message ID
	return "", fmt.Errorf("Not Implemented")
}
