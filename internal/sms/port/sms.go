package port

import (
	"context"
	"sms-dispatcher/internal/sms/domain"
)

type Repo interface {
	Create(ctx context.Context, SMS domain.SMS) (domain.SMSID, error)
	Update(ctx context.Context, SMS domain.SMS) error
	GetByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error)
}
