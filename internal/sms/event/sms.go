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

// publish update user balance event
type DebitBalanceEvent struct {
	Domain  Domain    `json:"domain"`
	EventID uuid.UUID `json:"event_id"`
	UserID  uint      `json:"user_id"`
	SMSID   uint      `json:"sms_id"`
	Amount  float64   `json:"amount"`
	Type    Type      `json:"type"`
}

// consumer event
type SMSUpdateEvent struct {
	Domain  Domain    `json:"domain"`
	EventID uuid.UUID `json:"event_id"`
	SMSID   uint      `json:"sms_id"`
	Status  Status    `json:"status"`
}

// refund event for finance
type RefundUserEvent struct {
	Domain  Domain    `json:"domain"`
	EventID uuid.UUID `json:"event_id"`
	UserID  uint      `json:"user_id"`
	SMSID   uint      `json:"sms_id"`
	Amount  float64   `json:"amount"`
	Type    Type      `json:"type"`
}
