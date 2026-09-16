package blogforger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	config "github.com/commitshark/notification-svc/internal"
	"github.com/commitshark/notification-svc/internal/domain/ports"
)

type BlogForgerTemplateRenderer struct {
	client *http.Client
	url    string
	bearer string
}

type EmailRenderResponse struct {
	Html      string `json:"html"`
	Subject   string `json:"subject"`
	PreHeader string `json:"preheader"`
}

func NewBlogForgerTemplateRenderer(config config.BlogForgerConfig) ports.TemplateRenderer {
	return &BlogForgerTemplateRenderer{
		url:    config.Url,
		bearer: config.Bearer,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (b *BlogForgerTemplateRenderer) Name() string {
	return "blogforger"
}

func (b *BlogForgerTemplateRenderer) Render(templateName string, subject string, data any, preHeader *string) (*ports.RenderResponse, error) {
	// template name is expected to be blogforger:{email_id}
	if !strings.HasPrefix(templateName, prefix) {
		return nil, fmt.Errorf("invalid template name %q: expected prefix %q", templateName, prefix)
	}

	emailID := strings.TrimPrefix(templateName, prefix)

	log.Println(fmt.Sprintf("Render email %s With data %v", emailID, data))

	url := fmt.Sprintf("/api/v1/emails/%s/build", emailID)

	payload := map[string]any{
		"variables": data,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", b.bearer)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed http call to email provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s renderer returned status %d", b.Name(), resp.StatusCode)
	}

	var responseBody EmailRenderResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return &ports.RenderResponse{
		Html:      responseBody.Html,
		Subject:   &responseBody.Subject,
		Preheader: &responseBody.PreHeader,
	}, nil
}
