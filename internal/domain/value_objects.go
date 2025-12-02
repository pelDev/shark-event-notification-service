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
	ID       string
	Email    *string
	Phone    *string
	DeviceID *string
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
	Title    string
	Body     *string
	Data     *map[string]interface{}
	HTML     *string
	Template *string
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
