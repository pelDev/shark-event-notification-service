package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/commitshark/notification-svc/internal/application/services"
	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/infrastructure/messaging"
)

// KafkaMessageHandler adapts Kafka messages to application service
type KafkaMessageHandler struct {
	service *services.NotificationService
	logger  *log.Logger
}

func NewKafkaMessageHandler(service *services.NotificationService) *KafkaMessageHandler {
	return &KafkaMessageHandler{
		service: service,
		logger:  log.New(log.Writer(), "[KafkaHandler] ", log.LstdFlags),
	}
}

// HandleMessage processes incoming Kafka messages
func (h *KafkaMessageHandler) HandleMessage(ctx context.Context, message []byte) error {
	var kafkaMsg messaging.KafkaNotificationMessage

	if err := json.Unmarshal(message, &kafkaMsg); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka message: %w", err)
	}

	// Validate message
	if err := kafkaMsg.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Convert to domain objects
	recipient, err := domain.NewRecipient(
		kafkaMsg.Recipient.ID,
		kafkaMsg.Recipient.Email,
		kafkaMsg.Recipient.Phone,
		kafkaMsg.Recipient.DeviceID,
	)
	if err != nil {
		return fmt.Errorf("invalid recipient: %w", err)
	}

	content, err := domain.NewContent(
		kafkaMsg.Content.Title,
		kafkaMsg.Content.Body,
		kafkaMsg.Content.Data,
	)
	if err != nil {
		return fmt.Errorf("invalid content: %w", err)
	}

	// Process through application service
	return h.service.ProcessNotification(
		ctx,
		kafkaMsg.ID,
		domain.NotificationType(kafkaMsg.Type),
		*recipient,
		*content,
		3,
	)
}
