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

func (r *smsRepo) GetByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
	var sms types.SMS
	query := r.db.WithContext(ctx).Model(&types.SMS{})

	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	tx := query.First(&sms)
	if tx.Error != nil {
		return nil, tx.Error
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
