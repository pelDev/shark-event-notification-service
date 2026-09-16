package blogforger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
)

type BlogForgerEmailProvider struct {
	URL      string
	Bearer   string
	SenderID string
	client   *http.Client
}

func NewBlogForgerEmailProvider(url, bearer, senderID string) *BlogForgerEmailProvider {
	return &BlogForgerEmailProvider{
		URL:      url,
		Bearer:   bearer,
		SenderID: senderID,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *BlogForgerEmailProvider) Name() string {
	return "blogforger-email-provider"
}

func (p *BlogForgerEmailProvider) Send(n *domain.Notification, isMarketing bool) (string, error) {
	if n.Recipient.Email == nil || *n.Recipient.Email == "" {
		return "", fmt.Errorf("email missing for SMS")
	}

	if n.Content.Template == nil || !strings.HasPrefix(*n.Content.Template, prefix) {
		return "", fmt.Errorf("invalid template name %s: expected prefix %q", *n.Content.Template, prefix)
	}

	if n.Content.Data == nil {
		return "", fmt.Errorf("invalid notification %s, data is empty", n.ID)
	}

	emailID := strings.TrimPrefix(*n.Content.Template, prefix)

	payload := map[string]any{
		"variables": *n.Content.Data,
		"to":        *n.Recipient.Email,
		"senderId":  p.SenderID,
	}

	log.Println(payload)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal notification: %w", err)
	}

	url := p.URL + fmt.Sprintf("/api/v1/emails/%s/send", emailID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", p.Bearer)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed http call to email provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("remote email provider returned status %d", resp.StatusCode)
	}

	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return fmt.Sprintf("email sent via blogforger provider: %v", responseBody), nil
}

func (p *BlogForgerEmailProvider) Supports(notificationType domain.NotificationType) bool {
	return notificationType == domain.EmailNotification
}
