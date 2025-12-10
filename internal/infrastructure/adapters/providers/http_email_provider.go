package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
)

type HttpEmailProvider struct {
	URL    string
	ApiKey string
	client *http.Client
}

func NewHTTPEmailProvider(url, apiKey string) *HttpEmailProvider {
	// Should inject http client instead??
	return &HttpEmailProvider{
		URL:    url,
		ApiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *HttpEmailProvider) Name() string {
	return "http-email-provider"
}

func (p *HttpEmailProvider) Send(n *domain.Notification) (string, error) {
	if n.Recipient.Email == nil || *n.Recipient.Email == "" {
		return "", fmt.Errorf("email missing for SMS")
	}

	payloadBytes, err := json.Marshal(n)
	if err != nil {
		return "", fmt.Errorf("failed to marshal notification: %w", err)
	}

	req, err := http.NewRequest("POST", p.URL, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", p.ApiKey)

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

	return fmt.Sprintf("email sent via http provider: %v", responseBody), nil
}

func (p *HttpEmailProvider) Supports(notificationType domain.NotificationType) bool {
	return notificationType == domain.EmailNotification
}
