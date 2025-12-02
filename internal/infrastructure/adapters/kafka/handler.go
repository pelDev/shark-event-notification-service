package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/commitshark/notification-svc/internal/application/services"
	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/events"
	"github.com/commitshark/notification-svc/internal/domain/ports"
	"github.com/commitshark/notification-svc/internal/infrastructure/messaging"
)

// KafkaMessageHandler adapts Kafka messages to application service
type KafkaMessageHandler struct {
	service        *services.NotificationService
	userDataSource ports.UserDataAdapter
	logger         *log.Logger
}

func NewKafkaMessageHandler(service *services.NotificationService, userAdapter ports.UserDataAdapter) *KafkaMessageHandler {
	return &KafkaMessageHandler{
		service:        service,
		userDataSource: userAdapter,
		logger:         log.New(log.Writer(), "[KafkaHandler] ", log.LstdFlags),
	}
}

// HandleMessage processes incoming Kafka messages
func (h *KafkaMessageHandler) HandleMessage(ctx context.Context, message []byte) error {
	var ev events.KafkaEvent

	if err := json.Unmarshal(message, &ev); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka message: %w", err)
	}

	if ev.EventType != "notification.requested" {
		h.logger.Println("Invalid event:", ev)
		return nil
	}

	// Validate message
	var payload messaging.KafkaNotificationMessagePayload
	if err := messaging.DecodeNotificationRequestPayload(&ev, &payload); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	var user, err = h.userDataSource.GetContactInfo(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("failed to get contact info message: %w", err)
	}

	// Convert to domain objects
	recipient, err := domain.NewRecipient(
		payload.UserID,
		&user.Email,
		user.Phone,
		user.DeviceID,
	)
	if err != nil {
		return fmt.Errorf("invalid recipient: %w", err)
	}

	content, err := domain.NewContent(
		payload.Subject,
		payload.Message,
		payload.Data,
		payload.HTML,
		payload.Template,
	)
	if err != nil {
		return fmt.Errorf("invalid content: %w", err)
	}

	// Process through application service
	return h.service.ProcessNotification(
		ctx,
		ev.ID,
		domain.NotificationType(payload.Channel),
		*recipient,
		*content,
		3,
	)
}
