package providers

import (
	"fmt"
	"mime"
	"net/smtp"
	"strconv"

	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/ports"
	domain_template "github.com/commitshark/notification-svc/internal/domain/templates"
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
		var emailData domain_template.EmailTemplateData
		err := domain_template.ParseTemplateData(*n.Content.Template, *n.Content.Data, &emailData)
		if err != nil {
			return "", err
		}

		html, err := p.renderer.Render(*n.Content.Template, subject, emailData, emailData.GetPreHeader())
		if err != nil {
			return "", err
		}

		fmt.Printf("[%s] Sending to %s: %s\n", p.Name(), email, subject)

		err = smtp.SendMail(
			p.smtpHost+":"+strconv.Itoa(p.smtpPort),
			p.smtpAuth,
			p.emailFrom,
			[]string{email},
			emailData.GetMessage(p.emailFrom, email, subject, html),
		)
		if err != nil {
			return "", err
		}

		fmt.Printf("Email sent successfully to %s", email)

		return fmt.Sprintf("Email sent successfully to %s", email), nil
	} else if n.Content.Body != nil && *n.Content.Body != "" {
		encodedSubject := mime.QEncoding.Encode("utf-8", subject)

		fmt.Println("Send plain mail:", encodedSubject)

		message := fmt.Sprintf(
			"From: \"Eventor\" <%s>\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
				"Reply-To: %s\r\n"+
				"\r\n"+
				"<!DOCTYPE html><html><body>%s</body></html>\r\n",
			p.emailFrom,
			email,
			encodedSubject,
			p.emailFrom,
			*n.Content.Body,
		)

		err := smtp.SendMail(
			p.smtpHost+":"+strconv.Itoa(p.smtpPort),
			p.smtpAuth,
			p.emailFrom,
			[]string{email},
			[]byte(message),
		)

		if err != nil {
			return "", err
		}

		fmt.Printf("Email sent successfully to %s\n", email)
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
