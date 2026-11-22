package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/ports"
	"github.com/commitshark/notification-svc/internal/domain/templates"
)

type NotificationService struct {
	repo                                 ports.NotificationRepository
	providers                            []ports.NotificationProvider
	blogForgerTicketPurchaseWithLocation *string
}

func NewNotificationService(
	repo ports.NotificationRepository,
	providers []ports.NotificationProvider,
	blogForgerTicketPurchaseWithLocation *string,
) *NotificationService {
	return &NotificationService{
		repo:                                 repo,
		providers:                            providers,
		blogForgerTicketPurchaseWithLocation: blogForgerTicketPurchaseWithLocation,
	}
}

// ProcessNotification processes incoming notification requests
func (s *NotificationService) ProcessNotification(
	ctx context.Context,
	id string,
	notificationType domain.NotificationType,
	recipient domain.Recipient,
	content domain.Content,
	maxRetries, isMarketing int,
) error {
	log.Printf("[ProcessNotification] Is notification marketing (1/0) = %d", isMarketing)

	// Create notification aggregate
	notification, err := domain.NewNotification(
		id,
		notificationType,
		recipient,
		content,
		maxRetries,
		isMarketing,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w, notificationType: %s", err, notificationType)
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

	needsBlogforger := notification.Content.Template != nil &&
		strings.HasPrefix(*notification.Content.Template, "blogforger")

	log.Printf("Need blogforger %v", needsBlogforger)

	// Find a provider that supports this notification type
	var provider ports.NotificationProvider
	for _, p := range s.providers {
		if !p.Supports(notification.Type) {
			continue
		}

		if needsBlogforger {
			if p.Name() == "blogforger-email-provider" {
				provider = p
				break
			}
			continue // don't fall back to other providers for blogforger templates
		}

		// non-blogforger template: pick any provider that isn't the blogforger one
		if !strings.HasPrefix(p.Name(), "blogforger") {
			provider = p
			break
		}
	}

	if provider == nil {
		return fmt.Errorf("no provider supports notification type %s", notification.Type)
	}

	log.Printf("[SendNotification] Send notification marketing (1/0) = %d Provider = %s", notification.IsMarketing, provider.Name())

	var providerResponse string
	var sendErr error

	log.Printf("Template = %s", *notification.Content.Template)

	// Attempt to send
	if notification.Content.Template != nil && *notification.Content.Template == "ticket-ready" {
		log.Println("handle send ticket ready")
		providerResponse, sendErr = s.handleSendTicketReadyEmail(notification, provider)
	} else {
		providerResponse, sendErr = provider.Send(notification, notification.IsMarketing == 1)
	}

	if sendErr != nil {
		notification.MarkAsFailed(sendErr.Error())
		if saveErr := s.repo.Save(ctx, notification); saveErr != nil {
			log.Printf("Failed to save failed notification: %v", saveErr)
		}
		return fmt.Errorf("failed to send notification: %w", sendErr)
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

func (s *NotificationService) handleSendTicketReadyEmail(
	n *domain.Notification,
	defaultProvider ports.NotificationProvider,
) (string, error) {
	if s.blogForgerTicketPurchaseWithLocation == nil || *s.blogForgerTicketPurchaseWithLocation == "" || n.Content.Data == nil {
		log.Println("fallback 1")
		return defaultProvider.Send(n, n.IsMarketing == 1)
	}

	raw, err := json.Marshal(*&n.Content.Data)
	if err != nil {
		return "", err
	}

	var notificationData templates.AttendeeTicketPurchaseEmailData
	if err := json.Unmarshal(raw, &notificationData); err != nil {
		return "", fmt.Errorf("failed to unmarshal ticket email data: %w", err)
	}

	log.Println(notificationData.IsRSVP)
	log.Println(notificationData.Location)

	if notificationData.IsRSVP {
		log.Println("fallback 2")
		return defaultProvider.Send(n, n.IsMarketing == 1)
	}

	blogforgerProvider := s.getBlogforgerProvider(n)

	if blogforgerProvider == nil {
		log.Println("fallback 3")
		return defaultProvider.Send(n, n.IsMarketing == 1)
	}

	nCopy := *n

	templateStr := fmt.Sprintf("blogforger:%s", *s.blogForgerTicketPurchaseWithLocation)
	log.Printf("Template string %s", templateStr)

	nCopy.Content.Template = &templateStr
	log.Println(nCopy.Content.Template)

	if notificationData.Location == nil || *notificationData.Location == "" {
		location := "Location would be communicated at most 2 days to the event"
		notificationData.Location = &location
		notificationDataMap, _ := structToMap(notificationData)
		nCopy.Content.Data = &notificationDataMap
	}

	return blogforgerProvider.Send(&nCopy, nCopy.IsMarketing == 1)
}

func (s *NotificationService) getBlogforgerProvider(n *domain.Notification) ports.NotificationProvider {
	var provider ports.NotificationProvider
	for _, p := range s.providers {
		if !p.Supports(n.Type) {
			continue
		}

		if p.Name() == "blogforger-email-provider" {
			provider = p
			break
		}
		continue // don't fall back to other providers for blogforger templates

	}
	return provider
}

func structToMap(v interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	return m, err
}
