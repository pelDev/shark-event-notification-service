package templates

import (
	"fmt"
	"strings"
	"time"
)

type AttendeeTicketPurchaseEmailData struct {
	TicketID   string  `json:"ticket_id"`
	QR         string  `json:"qr"`
	EventID    string  `json:"event_id"`
	EventTitle string  `json:"event_title"`
	TicketType string  `json:"ticket_type"`
	Date       string  `json:"date"`
	Amount     string  `json:"amount"`
	IsRSVP     bool    `json:"is_rsvp"`
	Location   *string `json:"location"`
	Admits     int     `json:"admits"`
}

func (tD *AttendeeTicketPurchaseEmailData) isEmailTemplateData() {}

func wrapBase64(s string, lineLen int) string {
	var buf strings.Builder
	for len(s) > 0 {
		chunk := lineLen
		if chunk > len(s) {
			chunk = len(s)
		}
		buf.WriteString(s[:chunk])
		buf.WriteString("\r\n")
		s = s[chunk:]
	}
	return buf.String()
}

func (tD *AttendeeTicketPurchaseEmailData) GetMessage(emailFrom, email, subject, html string) []byte {
	boundary := fmt.Sprintf("mixed-%d", time.Now().UnixNano())
	cid := "qr@local"

	// Strip data URL prefix if present
	qrBase64 := strings.TrimPrefix(tD.QR, "data:image/png;base64,")

	// ✅ Wrap at 76 chars per RFC 2045
	qrBase64Wrapped := wrapBase64(qrBase64, 76)

	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/related; boundary=\"%s\"\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: image/png\r\n"+
			"Content-ID: <%s>\r\n"+
			"Content-Disposition: inline; filename=\"qr.png\"\r\n"+ // ✅ Added
			"Content-Transfer-Encoding: base64\r\n"+
			"\r\n"+
			"%s\r\n"+ // ✅ Now wrapped
			"\r\n"+
			"--%s--\r\n",
		emailFrom,
		email,
		subject,
		boundary,
		boundary,
		html,
		boundary,
		cid,
		qrBase64Wrapped, // ✅ Use wrapped version
		boundary,
	)

	return []byte(message)
}

func (tD *AttendeeTicketPurchaseEmailData) GetPreHeader() *string {
	return nil
}
