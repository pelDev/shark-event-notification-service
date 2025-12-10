package domain

import (
	"errors"
)

type UserContactInfo struct {
	Email    string
	Phone    *string
	DeviceID *string
}

type Recipient struct {
	ID       string  `json:"id"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	DeviceID *string `json:"device_id,omitempty"`
}

func NewRecipient(id string, email, phone, deviceId *string) (*Recipient, error) {
	if id == "" {
		return nil, errors.New("recipient id cannot be empty")
	}
	return &Recipient{
		ID:       id,
		Email:    email,
		Phone:    phone,
		DeviceID: deviceId,
	}, nil
}

type Content struct {
	Title    string                  `json:"title"`
	Body     *string                 `json:"body,omitempty"`
	Data     *map[string]interface{} `json:"data,omitempty"`
	HTML     *string                 `json:"html,omitempty"`
	Template *string                 `json:"template,omitempty"`
}

func NewContent(title string, body *string, data *map[string]interface{}, html, template *string) (*Content, error) {
	if title == "" {
		return nil, errors.New("title and body cannot be empty")
	}

	if (body == nil || *body == "") && (data == nil) {
		return nil, errors.New("data and body cannot be empty")
	}

	return &Content{
		Title:    title,
		Body:     body,
		Data:     data,
		HTML:     html,
		Template: template,
	}, nil
}
