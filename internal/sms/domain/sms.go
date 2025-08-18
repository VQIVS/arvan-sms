package domain

import "time"

type Status string

const (
	Pending   Status = "pending"
	Delivered Status = "delivered"
	Failed    Status = "failed"
)

type SMSFilter struct {
	ID     SMSID
	Status Status
}

type SMSID uint

type SMS struct {
	ID        SMSID     `json:"id"`
	Recipient string    `json:"recipient"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
