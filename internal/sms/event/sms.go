package event

import "github.com/google/uuid"

type Type string
type Domain string
type Status string

const (
	SMSDebitEvent  Type   = "sms_debit"
	SMSCreditEvent Type   = "sms_credit"
	SMS            Domain = "sms"
	Finance        Domain = "finance"
	StatusSuccess  Status = "success"
	StatusFailed   Status = "failed"
)

// Publisher Event
type UserBalanceEvent struct {
	Domain  Domain    `json:"domain"`
	EventID uuid.UUID `json:"event_id"`
	UserID  uint      `json:"user_id"`
	SMSID   uint      `json:"sms_id"`
	Amount  float64   `json:"amount"`
	Type    Type      `json:"type"`
}

// Consumer Event
type SMSUpdateEvent struct {
	Domain  Domain    `json:"domain"`
	EventID uuid.UUID `json:"event_id"`
	SMSID   uint      `json:"sms_id"`
	Status  Status    `json:"status"`
}
