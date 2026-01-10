package models

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	OrgID        string
	CreatedAt    time.Time
}

type Organization struct {
	ID          string
	Name        string
	OwnerUserID string
	CreatedAt   time.Time
}

type Client struct {
	ID        string
	OrgID     string
	Name      string
	Email     string
	Company   string
	Phone     string
	Notes     string
	CreatedAt time.Time
}

type Template struct {
	ID        string
	OrgID     string
	Name      string
	Subject   string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Invoice struct {
	ID         string
	OrgID      string
	ClientID   string
	TemplateID string
	Number     string
	AmountCents int64
	Currency   string
	DueDate    time.Time
	Status     string
	Notes      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Reminder struct {
	ID           string
	OrgID        string
	InvoiceID    string
	TemplateID   string
	ScheduledFor time.Time
	SentAt       *time.Time
	Status       string
	CreatedAt    time.Time
}

type OutboxEmail struct {
	ID         string
	OrgID      string
	ReminderID string
	ToEmail    string
	Subject    string
	Body       string
	CreatedAt  time.Time
}
