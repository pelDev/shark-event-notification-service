package templates

import (
	"encoding/json"
	"fmt"
)

type EmailTemplateData interface {
	isEmailTemplateData()
	GetMessage(emailFrom, email, subject, html string) []byte
}

func ParseTemplateData(templateName string, data map[string]interface{}, out *EmailTemplateData) error {
	fmt.Printf("[ParseTemplateData] templateName: %s, data: %v", templateName, data)

	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	switch templateName {
	case "ticket-ready":
		var result AttendeeTicketPurchaseEmailData
		if err := json.Unmarshal(raw, &result); err != nil {
			return fmt.Errorf("failed to unmarshal ticket email data: %w", err)
		}
		*out = &result
		return nil

	case "otp":
		var result OtpData
		if err := json.Unmarshal(raw, &result); err != nil {
			return fmt.Errorf("failed to unmarshal otp email data: %w", err)
		}
		*out = &result
		return nil

	default:
		return fmt.Errorf("unknown template name: %s", templateName)
	}
}
