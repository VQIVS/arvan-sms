package dto

import "time"

type SendSMSRequest struct {
	Content  string `json:"content" validate:"required,max=160"`
	Receiver string `json:"receiver" validate:"required,e164"` // E.164 format phone number
	UserID   string `json:"user_id" validate:"required,uuid"`
}

type SendSMSResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message,omitempty"`
}

type GetSMSResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Content     string     `json:"content"`
	Receiver    string     `json:"receiver"`
	Provider    string     `json:"provider,omitempty"`
	Status      string     `json:"status"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	FailureCode string     `json:"failure_code,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
