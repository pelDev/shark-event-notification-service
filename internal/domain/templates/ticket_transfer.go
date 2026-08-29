package templates

import (
	"fmt"
)

// ==================== RECIPIENT (Ticket Transfer In) ====================

type TicketTransferInData struct {
	TransferID    string `json:"transfer_id"`
	TicketID      string `json:"ticket_id"`
	TicketType    string `json:"ticket_type"`
	EventTitle    string `json:"event_title"`
	FromUser      string `json:"from_user"`
	FromUserEmail string `json:"from_user_email"`
	TransferDate  string `json:"transfer_date"`
	TransferTime  string `json:"transfer_time"`
	Amount        string `json:"amount"`
	QR            string `json:"qr"`
	Admits        int    `json:"admits"`
	Message       string `json:"message"`
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
		e.FromUser,
		e.EventTitle)
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
	TransferID   string  `json:"transfer_id"`
	TicketID     string  `json:"ticket_id"`
	TicketType   string  `json:"ticket_type"`
	EventTitle   string  `json:"event_title"`
	ToUser       string  `json:"to_user"`
	ToUserEmail  *string `json:"to_user_email"`
	TransferDate string  `json:"transfer_date"`
	TransferTime string  `json:"transfer_time"`
	Amount       string  `json:"amount"`
	Admits       int     `json:"admits"`
	Message      string  `json:"message"`
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
		e.EventTitle,
		e.ToUser)
	return &preHeader
}
