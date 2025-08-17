package storage

import (
	"context"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/port"
	"sms-dispatcher/pkg/adapters/storage/mapper"
	"sms-dispatcher/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

type smsRepo struct {
	db *gorm.DB
}

func NewSMSRepo(db *gorm.DB) port.Repo {
	return &smsRepo{
		db: db,
	}
}

func (r *smsRepo) Create(ctx context.Context, SMS domain.SMS) (domain.SMSID, error) {
	sms := mapper.SMSDomain2Storage(SMS)
	tx := r.db.WithContext(ctx).Create(sms)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return domain.SMSID(sms.ID), nil
}

func (r *smsRepo) GetByID(ctx context.Context, id domain.SMSID) (*domain.SMS, error) {
	var sms types.SMS
	if err := r.db.WithContext(ctx).First(&sms, id).Error; err != nil {
		return nil, err
	}
	return mapper.SMSStorage2Domain(&sms), nil
}

func (r *smsRepo) Update(ctx context.Context, SMS domain.SMS) error {
	sms := mapper.SMSDomain2Storage(SMS)
	tx := r.db.WithContext(ctx).Save(sms)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
