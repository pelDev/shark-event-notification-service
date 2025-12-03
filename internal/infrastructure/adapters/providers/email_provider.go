package providers

import (
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/ports"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/templates"
)

type EmailProvider struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	emailFrom    string
	smtpAuth     smtp.Auth
	renderer     ports.TemplateRenderer
}

func NewEmailProvider(host string, port int, username, password, from string, renderer ports.TemplateRenderer, auth smtp.Auth) *EmailProvider {
	return &EmailProvider{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		emailFrom:    from,
		renderer:     renderer,
		smtpAuth:     auth,
	}
}

func (p *EmailProvider) Name() string {
	return "Custom Smtp"
}

func (p *EmailProvider) Send(n *domain.Notification) (string, error) {
	if n.Recipient.Email == nil || *n.Recipient.Email == "" {
		return "", fmt.Errorf("email missing for SMS")
	}

	email := *n.Recipient.Email
	subject := n.Content.Title

	if n.Content.Template != nil && *n.Content.Template != "" && n.Content.Data != nil {
		var emailData templates.EmailTemplateData
		err := templates.ParseTemplateData(*n.Content.Template, *n.Content.Data, &emailData)
		if err != nil {
			return "", err
		}

		html, err := p.renderer.Render(*n.Content.Template, subject, emailData)
		if err != nil {
			return "", err
		}

		fmt.Printf("[%s] Sending to %s: %s\n", p.Name(), email, subject)

		message := []byte(fmt.Sprintf(
			"From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s",
			p.emailFrom, email, subject, html,
		))

		err = smtp.SendMail(p.smtpHost+":"+strconv.Itoa(p.smtpPort), p.smtpAuth, p.emailFrom, []string{email}, message)
		if err != nil {
			return "", err
		}

		fmt.Printf("Email sent successfully to %s", email)

		return fmt.Sprintf("Email sent successfully to %s", email), nil
	}

	// html, err := p.renderer.Render("email_ticket.html", templateData)
	// if err != nil {
	// 	return "", fmt.Errorf("template rendering failed: %w", err)
	// }

	// Simulate email sending

	// In production, integrate with email service (SendGrid, AWS SES, etc.)

	// body := notification.Content.Body

	fmt.Printf("[%s] Send to %s: %s\n", p.Name(), email, subject)

	// return fmt.Sprintf("Email sent to %s with subject: %s", email, subject), nil
	return "", fmt.Errorf("Failed to send email")
}

func (p *EmailProvider) Supports(notificationType domain.NotificationType) bool {
	return notificationType == domain.EmailNotification
}
