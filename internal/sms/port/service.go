package port

import (
	"context"
	"sms-dispatcher/internal/sms/domain"
)

type Service interface {
	SendSMS(ctx context.Context, recipient string, message string) (string, error)
	GetSMSStatus(ctx context.Context, smsID string) (string, error)
	ListSMS(ctx context.Context, status string) ([]domain.SMS, error)
	/**
	PublishSMSDebit is a method to publish an SMS debit event to finance
	event {
	domain: "sms",
	ID: smsID,
	amount: 1.0, // Assuming a fixed cost per SMS
	}
	finance acks then updates the SMS status to "delivered" or "failed" based on the ack.
	**/
	PublishSMSDebit(ctx context.Context, smsID string, amount float64) error
	/**
	ConsumeSMSUpdate is a method to consume SMS update events.
	It is used to update the SMS status based on events from the finance service.
	event {
		domain: "sms",
		ID: smsID,
		status: "delivered" | "failed"
		}
	It takes the SMS ID as an argument and updates the status accordingly.
	It is typically used in a message queue or event-driven architecture.
	**/
	ConsumeSMSUpdate(ctx context.Context, smsID string) error
}
