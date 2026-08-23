package templates

import (
	"fmt"
)

// ==================== RECIPIENT (Ticket Transfer In) ====================

type TicketTransferInEvent struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type TicketTransferInTicket struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Admits     int    `json:"admits"`
	Amount     string `json:"amount"`
	QR         string `json:"qr"`
	TransferID string `json:"transfer_id"`
}

type TicketTransferInFromUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TicketTransferInData struct {
	Event        TicketTransferInEvent    `json:"event"`
	Ticket       TicketTransferInTicket   `json:"ticket"`
	FromUser     TicketTransferInFromUser `json:"from_user"`
	TransferDate string                   `json:"transfer_date"`
	TransferTime string                   `json:"transfer_time"`
	Message      string                   `json:"message"`
}

func (e *TicketTransferInData) isEmailTemplateData() {}

func (e *TicketTransferInData) GetMessage(emailFrom, email, subject, html string) []byte {
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		emailFrom,
		email,
		subject,
		html,
	)

	return []byte(message)
}

func (e *TicketTransferInData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("🎟️ %s has transferred a ticket to you for '%s'!",
		e.FromUser.Name,
		e.Event.Name)
	return &preHeader
}

// ==================== SENDER (Ticket Transfer Out) ====================

type TicketTransferOutEvent struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type TicketTransferOutTicket struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Admits     int    `json:"admits"`
	Amount     string `json:"amount"`
	TransferID string `json:"transfer_id"`
}

type TicketTransferOutToUser struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

type TicketTransferOutData struct {
	Event        TicketTransferOutEvent  `json:"event"`
	Ticket       TicketTransferOutTicket `json:"ticket"`
	ToUser       TicketTransferOutToUser `json:"to_user"`
	TransferDate string                  `json:"transfer_date"`
	TransferTime string                  `json:"transfer_time"`
	Message      string                  `json:"message"`
}

func (e *TicketTransferOutData) isEmailTemplateData() {}

func (e *TicketTransferOutData) GetMessage(emailFrom, email, subject, html string) []byte {
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		emailFrom,
		email,
		subject,
		html,
	)

	return []byte(message)
}

func (e *TicketTransferOutData) GetPreHeader() *string {
	preHeader := fmt.Sprintf("✅ You successfully transferred a ticket for '%s' to %s!",
		e.Event.Name,
		e.ToUser.Name)
	return &preHeader
}
