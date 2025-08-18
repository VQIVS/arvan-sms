package port

import (
	"context"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
)

type Service interface {
	CreateSMS(ctx context.Context, recipient string, message string) (domain.SMSID, error)
	GetSMSByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error)
	//TODO: Implement ListSMS method

	/*
		PublishSMSDebit is a method to publish an SMS debit event to finance
		event {
		domain: "sms",
		ID: smsID,
		amount: 1.0, // Assuming a fixed cost per SMS
		}
		finance acks then updates the SMS status to "delivered" or "failed" based on the ack.
	*/
	UserBalanceUpdate(ctx context.Context, user event.UserBalanceEvent) error
	/*
		ConsumeSMSUpdate is a method to consume SMS update events.
		It is used to update the SMS status based on events from the finance service.
		event {
			domain: "sms",
			ID: smsID,
			status: "delivered" | "failed"
			}
		It takes the SMS ID as an argument and updates the status accordingly.
		It is typically used in a message queue or event-driven architecture.
	*/
	UpdateSMSStatus(ctx context.Context, body []byte) error
}
