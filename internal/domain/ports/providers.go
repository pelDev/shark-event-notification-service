package ports

import "github.com/commitshark/notification-svc/internal/domain"

type NotificationProvider interface {
	Send(notification *domain.Notification) (string, error)
	Supports(notificationType domain.NotificationType) bool
	Name() string
}

type TemplateRenderer interface {
	Render(templateName, subject string, data any, preHeader *string) (string, error)
}
