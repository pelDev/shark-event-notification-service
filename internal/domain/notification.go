package domain

import (
	"errors"
	"time"
)

type Notification struct {
	ID               string
	Type             NotificationType
	Recipient        Recipient
	Content          Content
	Status           NotificationStatus
	ProviderResponse string
	CreatedAt        time.Time
	SentAt           *time.Time
	RetryCount       int
	MaxRetries       int
	Version          int // For optimistic locking
}

// Business rules
func (n *Notification) CanBeSent() bool {
	return n.Status == StatusPending ||
		(n.Status == StatusFailed && n.RetryCount < n.MaxRetries)
}

func (n *Notification) MarkAsSent(providerResponse string) error {
	if n.Status == StatusSent || n.Status == StatusDelivered {
		return errors.New("notification already sent or delivered")
	}

	now := time.Now()
	n.Status = StatusSent
	n.SentAt = &now
	n.ProviderResponse = providerResponse
	n.Version++

	return nil
}

func (n *Notification) MarkAsFailed(providerResponse string) {
	n.Status = StatusFailed
	n.RetryCount++
	n.ProviderResponse = providerResponse
	n.Version++
}

func (n *Notification) MarkAsDelivered() error {
	if n.Status != StatusSent {
		return errors.New("only sent notifications can be marked as delivered")
	}

	n.Status = StatusDelivered
	n.Version++

	// Could add NotificationDelivered event here
	return nil
}

// Factory method with validation
func NewNotification(
	id string,
	notificationType NotificationType,
	recipient Recipient,
	content Content,
	maxRetries int,
) (*Notification, error) {

	if id == "" {
		return nil, errors.New("notification id cannot be empty")
	}

	if !isValidNotificationType(notificationType) {
		return nil, errors.New("invalid notification type")
	}

	if maxRetries < 0 {
		maxRetries = 3 // default
	}

	return &Notification{
		ID:         id,
		Type:       notificationType,
		Recipient:  recipient,
		Content:    content,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
		MaxRetries: maxRetries,
		RetryCount: 0,
		Version:    1,
	}, nil
}

func isValidNotificationType(t NotificationType) bool {
	switch t {
	case EmailNotification, SMSNotification, PushNotification:
		return true
	default:
		return false
	}
}
