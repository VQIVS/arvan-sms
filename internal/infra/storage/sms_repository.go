package storage

import (
	"context"
	"sms/internal/domain/sms"
	"sms/internal/infra/storage/mapper"
	"sms/internal/infra/storage/types"

	"gorm.io/gorm"
)

type SMSRepository struct {
	Db *gorm.DB
}

func NewSMSRepository(db *gorm.DB) sms.Repo {
	return &SMSRepository{
		Db: db,
	}
}

func (r *SMSRepository) WithTx(tx *gorm.DB) sms.Repo {
	return &SMSRepository{
		Db: tx,
	}
}

func (r *SMSRepository) GetByFilter(ctx context.Context, filter sms.Filter) (*sms.SMSMessage, error) {
	var sms types.SMS
	query := r.Db.WithContext(ctx)
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if err := query.First(&sms).Error; err != nil {
		return nil, err
	}
	return mapper.TODomain(sms), nil
}

func (r *SMSRepository) Create(ctx context.Context, message *sms.SMSMessage) error {
	model := mapper.TOStorage(*message)
	return r.Db.WithContext(ctx).Create(&model).Error
}

func (r *SMSRepository) Update(ctx context.Context, ID string, message *sms.SMSMessage) error {
	var model types.SMS
	return r.Db.
		WithContext(ctx).
		Model(&model).
		Where("id = ?", ID).
		Updates(mapper.TOStorage(*message)).Error
}
