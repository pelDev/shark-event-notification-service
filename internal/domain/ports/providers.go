package ports

import "github.com/commitshark/notification-svc/internal/domain"

type NotificationProvider interface {
	Send(notification *domain.Notification, isMarketing bool) (string, error)
	Supports(notificationType domain.NotificationType) bool
	Name() string
}

type RenderResponse struct {
	Html      string
	Subject   *string
	Preheader *string
}

type TemplateRenderer interface {
	Name() string
	Render(templateName, subject string, data any, preHeader *string) (*RenderResponse, error)
}
