package sms

import (
	"context"
	"time"
)

type EventType string

const (
	EventTypeDebit     EventType = "Debit"
	EventTypeRefund    EventType = "Refund"
	EventTypeDelivered EventType = "Delivered"
	EventTypeFailed    EventType = "Failed"
)

type Publisher interface {
	PublishEvent(ctx context.Context, event SMSEvent) error
}

type SMSEvent interface {
	EventType() EventType
	AggregateID() string
}

type SMSDelivered struct {
	SMSID     string    `json:"sms_id"`
	Provider  string    `json:"provider"`
	TimeStamp time.Time `json:"timestamp"`
}

type SMSFailed struct {
	SMSID       string    `json:"sms_id"`
	Provider    string    `json:"provider"`
	FailureCode string    `json:"failure_code"`
	Reason      string    `json:"reason"`
	TimeStamp   time.Time `json:"timestamp"`
}

type DebitUserBalance struct {
	UserID    string    `json:"user_id"`
	SMSID     string    `json:"sms_id"`
	Amount    float64   `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}

func (e SMSDelivered) EventType() EventType {
	return EventTypeDelivered
}

func (e SMSFailed) EventType() EventType {
	return EventTypeFailed
}

func (e DebitUserBalance) EventType() EventType {
	return EventTypeDebit
}

func (e SMSDelivered) AggregateID() string {
	return e.SMSID
}

func (e SMSFailed) AggregateID() string {
	return e.SMSID
}

func (e DebitUserBalance) AggregateID() string {
	return e.SMSID
}
