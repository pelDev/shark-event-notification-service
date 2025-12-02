package domain

import (
	"errors"
)

type Recipient struct {
	ID    string
	Email string
	Phone string
}

func NewRecipient(id, email, phone string) (*Recipient, error) {
	if id == "" {
		return nil, errors.New("recipient id cannot be empty")
	}
	return &Recipient{
		ID:    id,
		Email: email,
		Phone: phone,
	}, nil
}

type Content struct {
	Title string
	Body  string
	Data  map[string]interface{}
}

func NewContent(title, body string, data map[string]interface{}) (*Content, error) {
	if title == "" || body == "" {
		return nil, errors.New("title and body cannot be empty")
	}
	return &Content{
		Title: title,
		Body:  body,
		Data:  data,
	}, nil
}
