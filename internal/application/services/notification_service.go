package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/ports"
)

type NotificationService struct {
	repo      ports.NotificationRepository
	providers []ports.NotificationProvider
}

func NewNotificationService(
	repo ports.NotificationRepository,
	providers []ports.NotificationProvider,
) *NotificationService {
	return &NotificationService{
		repo:      repo,
		providers: providers,
	}
}

// ProcessNotification processes incoming notification requests
func (s *NotificationService) ProcessNotification(
	ctx context.Context,
	id string,
	notificationType domain.NotificationType,
	recipient domain.Recipient,
	content domain.Content,
	maxRetries int,
) error {

	// Create notification aggregate
	notification, err := domain.NewNotification(
		id,
		notificationType,
		recipient,
		content,
		maxRetries,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Save to repository
	if err := s.repo.Save(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notification (TODO: could be async)
	go func() {
		if err := s.SendNotification(context.Background(), notification.ID); err != nil {
			log.Printf("Failed to send notification %s: %v", notification.ID, err)
		}
	}()

	return nil
}

// SendNotification attempts to send a notification
func (s *NotificationService) SendNotification(ctx context.Context, notificationID string) error {
	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("notification not found: %w", err)
	}

	if !notification.CanBeSent() {
		return fmt.Errorf("notification %s cannot be sent", notification.ID)
	}

	// Find a provider that supports this notification type
	var provider ports.NotificationProvider
	for _, p := range s.providers {
		if p.Supports(notification.Type) {
			provider = p
			break
		}
	}

	if provider == nil {
		return fmt.Errorf("no provider supports notification type %s", notification.Type)
	}

	// Attempt to send
	providerResponse, err := provider.Send(notification)
	if err != nil {
		notification.MarkAsFailed(err.Error())
		if saveErr := s.repo.Save(ctx, notification); saveErr != nil {
			log.Printf("Failed to save failed notification: %v", saveErr)
		}
		return fmt.Errorf("failed to send notification: %w", err)
	}

	// Mark as sent
	if err := notification.MarkAsSent(providerResponse); err != nil {
		return err
	}

	// Save changes
	if err := s.repo.Save(ctx, notification); err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

// RetryFailedNotifications retries failed notifications with exponential backoff
func (s *NotificationService) RetryFailedNotifications(ctx context.Context, batchSize int) error {
	pending, err := s.repo.FindPending(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to find pending notifications: %w", err)
	}

	for _, notification := range pending {
		// Calculate backoff: 5s * 2^retryCount
		backoff := time.Duration(5*(1<<notification.RetryCount)) * time.Second

		go func(n *domain.Notification, delay time.Duration) {
			time.Sleep(delay)
			if err := s.SendNotification(context.Background(), n.ID); err != nil {
				log.Printf("Retry failed for notification %s: %v", n.ID, err)
			}
		}(notification, backoff)
	}

	return nil
}
