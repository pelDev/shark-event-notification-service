package domain

type NotificationType string

const (
	EmailNotification NotificationType = "EMAIL"
	SMSNotification   NotificationType = "SMS"
	PushNotification  NotificationType = "PUSH"
	InAppNotification NotificationType = "IN_APP"
)

type NotificationStatus string

const (
	StatusPending   NotificationStatus = "PENDING"
	StatusSent      NotificationStatus = "SENT"
	StatusFailed    NotificationStatus = "FAILED"
	StatusDelivered NotificationStatus = "DELIVERED"
)
