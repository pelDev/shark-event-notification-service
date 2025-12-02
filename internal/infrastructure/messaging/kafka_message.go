package messaging

import "fmt"

// Infrastructure DTO for Kafka messages
type KafkaNotificationMessage struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Recipient struct {
		ID       string  `json:"id"`
		Email    *string `json:"email"`
		Phone    *string `json:"phone"`
		DeviceID *string `json:"device_id"`
	} `json:"recipient"`
	Content struct {
		Title string                  `json:"title"`
		Body  *string                 `json:"body"`
		Data  *map[string]interface{} `json:"data"`
	} `json:"content"`
}

func (m *KafkaNotificationMessage) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("id is required")
	}
	if m.Type == "" {
		return fmt.Errorf("type is required")
	}
	if m.Recipient.ID == "" {
		return fmt.Errorf("recipient.id is required")
	}
	if m.Content.Title == "" {
		return fmt.Errorf("content.title is required")
	}
	if (m.Content.Body == nil || *m.Content.Body == "") && (m.Content.Data == nil) {
		return fmt.Errorf("content.body and content.data cannot be empty")
	}
	return nil
}
