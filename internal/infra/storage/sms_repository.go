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

func (r *SMSRepository) GetByID(ctx context.Context, ID string) (*sms.SMSMessage, error) {
	var sms types.SMS
	if err := r.Db.WithContext(ctx).Where("id = ?", ID).First(&sms).Error; err != nil {
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
