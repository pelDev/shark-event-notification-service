package applicationdto

import (
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
)

// ListNotificationsRequest represents the query parameters for listing notifications
type ListNotificationsRequest struct {
	Page        int    `json:"-"` // from query param
	PageSize    int    `json:"-"` // from query param
	Status      string `json:"-"` // from query param, optional filter
	Type        string `json:"-"` // from query param, optional filter
	IsMarketing *bool  `json:"-"`
}

type NotificationDto struct {
	ID               string                    `json:"id"`
	Type             domain.NotificationType   `json:"type"`
	Recipient        domain.Recipient          `json:"recipient"`
	ContentTitle     string                    `json:"content_title"`
	Status           domain.NotificationStatus `json:"status"`
	ProviderResponse string                    `json:"provider_response"`
	CreatedAt        time.Time                 `json:"created_at"`
	SentAt           *time.Time                `json:"sent_at,omitempty"`
	RetryCount       int                       `json:"retry_count"`
	MaxRetries       int                       `json:"max_retries"`
	Version          int                       `json:"version"`
	IsMarketing      int                       `json:"is_marketing"`
}

func ToNotificationDtos(notifications []*domain.Notification) []*NotificationDto {
	dtos := make([]*NotificationDto, 0, len(notifications))

	for _, n := range notifications {
		dtos = append(dtos, &NotificationDto{
			ID:               n.ID,
			Type:             n.Type,
			Recipient:        n.Recipient,
			ContentTitle:     n.Content.Title,
			Status:           n.Status,
			ProviderResponse: n.ProviderResponse,
			CreatedAt:        n.CreatedAt,
			SentAt:           n.SentAt,
			RetryCount:       n.RetryCount,
			MaxRetries:       n.MaxRetries,
			Version:          n.Version,
			IsMarketing:      n.IsMarketing,
		})
	}

	return dtos
}
