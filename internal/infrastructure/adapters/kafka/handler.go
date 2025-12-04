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
// This handler assumes the message is a DomainEvent with NotificationMessagePayload
func (h *KafkaMessageHandler) HandleMessage(ctx context.Context, ev events.DomainEvent) error {
	var payload events.NotificationMessagePayload
	if err := json.Unmarshal(ev.Payload, &payload); err != nil {
		return fmt.Errorf("notification[%s]: invalid payload: %w", ev.ID, err)
	}

	user, err := h.userDataSource.GetContactInfo(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("notification[%s]: getContactInfo failed: %w", ev.ID, err)
	}

	// Convert to domain objects
	recipient, err := domain.NewRecipient(
		payload.UserID,
		&user.Email,
		user.Phone,
		user.DeviceID,
	)
	if err != nil {
		return fmt.Errorf("notification[%s]: invalid content: %w", ev.ID, err)
	}

	content, err := domain.NewContent(
		payload.Subject,
		payload.Message,
		payload.Data,
		payload.HTML,
		payload.Template,
	)
	if err != nil {
		return fmt.Errorf("notification[%s]: invalid content: %w", ev.ID, err)
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
