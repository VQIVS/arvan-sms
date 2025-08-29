package sms

import "time"

type SMSEvent interface {
	EventType() string
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

func (e SMSDelivered) EventType() string {
	return "SMSDelivered"
}

func (e SMSFailed) EventType() string {
	return "SMSFailed"
}

func (e SMSDelivered) AggregateID() string {
	return e.SMSID
}

func (e SMSFailed) AggregateID() string {
	return e.SMSID
}
