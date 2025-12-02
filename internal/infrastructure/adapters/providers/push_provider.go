package providers

import (
	"fmt"

	"github.com/commitshark/notification-svc/internal/domain"
)

type PushProvider struct {
}

func NewPushProvider() *PushProvider {
	return &PushProvider{}
}

func (p *PushProvider) Name() string {
	return "Sample Push Provider"
}

func (p *PushProvider) Supports(notificationType domain.NotificationType) bool {
	return notificationType == domain.PushNotification
}

func (p *PushProvider) Send(n *domain.Notification) (string, error) {
	if n.Recipient.DeviceID == nil || *n.Recipient.DeviceID == "" {
		return "", fmt.Errorf("device id missing for SMS")
	}

	if n.Content.Body == nil || *n.Content.Body == "" {
		return "", fmt.Errorf("message content for SMS")
	}

	// Example; replace with actual API call
	fmt.Printf("[SMS] Sending to %s: %s\n", *n.Recipient.DeviceID, *n.Content.Body)

	// Return a fake provider message ID
	return "", fmt.Errorf("Not Implemented")
}
