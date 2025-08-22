package port

import (
	"context"
	"sms-dispatcher/internal/sms/domain"

	"gorm.io/gorm"
)

type Repo interface {
	Create(ctx context.Context, SMS domain.SMS) (domain.SMSID, error)
	Update(ctx context.Context, SMS domain.SMS) error
	GetByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error)
	// Aplying Tx to the repository behavior
	BeginTx(ctx context.Context) (*gorm.DB, error)
}
