package sms

import (
	"context"
	"time"
)

type EventType string

const (
	EventTypeBillingRequested EventType = "BillingRequested"
	EventTypeBillingCompleted EventType = "BillingCompleted"
	EventTypeBillingRefunded  EventType = "BillingRefunded"
)

type EventPublisher interface {
	PublishEvent(ctx context.Context, event DomainEvent) error
}

type DomainEvent interface {
	EventType() EventType

	AggregateID() string

	Timestamp() time.Time
}

type RequestSMSBilling struct {
	UserID    string    `json:"user_id"`
	SMSID     string    `json:"sms_id"`
	Amount    int64     `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}

func (e RequestSMSBilling) EventType() EventType {
	return EventTypeBillingRequested
}

func (e RequestSMSBilling) AggregateID() string {
	return e.SMSID
}

func (e RequestSMSBilling) Timestamp() time.Time {
	return e.TimeStamp
}

type SMSBillingCompleted struct {
	UserID        string    `json:"user_id"`
	SMSID         string    `json:"sms_id"`
	Amount        int64     `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	TimeStamp     time.Time `json:"timestamp"`
}

func (e SMSBillingCompleted) EventType() EventType {
	return EventTypeBillingCompleted
}

func (e SMSBillingCompleted) AggregateID() string {
	return e.SMSID
}

func (e SMSBillingCompleted) Timestamp() time.Time {
	return e.TimeStamp
}

type RequestBillingRefund struct {
	TransactionID string    `json:"transaction_id"`
	TimeStamp     time.Time `json:"timestamp"`
}

func (e RequestBillingRefund) EventType() EventType {
	return EventTypeBillingRefunded
}

func (e RequestBillingRefund) AggregateID() string {
	return e.TransactionID
}

func (e RequestBillingRefund) Timestamp() time.Time {
	return e.TimeStamp
}
