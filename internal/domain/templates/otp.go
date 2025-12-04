package templates

type OtpData struct {
	OtpCode  string `json:"otp_code"`
	ValidFor string `json:"valid_for"` // e.g. "5 minutes"
}
