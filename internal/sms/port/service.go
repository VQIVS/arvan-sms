package port

import (
	"context"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
)

type Service interface {
	CreateSMS(ctx context.Context, recipient string, message string) (domain.SMSID, error)
	GetSMSByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error)
	UserBalanceUpdate(ctx context.Context, user event.UserBalanceEvent) error
	UpdateSMSStatus(ctx context.Context, body []byte) error
}
