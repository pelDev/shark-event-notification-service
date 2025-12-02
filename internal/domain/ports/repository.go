package ports

import (
	"context"

	"github.com/commitshark/notification-svc/internal/domain"
)

// Domain defines what it needs - Repository port
type NotificationRepository interface {
	Save(ctx context.Context, notification *domain.Notification) error
	FindByID(ctx context.Context, id string) (*domain.Notification, error)
	FindPending(ctx context.Context, limit int) ([]*domain.Notification, error)
	UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus, providerResponse string) error
	IncrementRetryCount(ctx context.Context, id string) error
}

type UserDataAdapter interface {
	GetContactInfo(ctx context.Context, userID string) (*domain.UserContactInfo, error)
}
