package templates

type TicketSoldData struct {
	EventTitle string `json:"event_title"`
	TicketType string `json:"ticket_type"`
	BuyerName  string `json:"buyer_name"`
	BuyerEmail string `json:"buyer_email"`
	Amount     string `json:"amount"`
	TicketID   string `json:"ticket_id"`
}

// TODO: Implement EmailTemplateData interface methods
