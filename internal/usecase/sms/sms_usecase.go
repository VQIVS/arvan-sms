package sms

import (
	"context"
	"sms/internal/domain/sms"
	"time"

	"gorm.io/gorm"
)

type UseCase struct {
	smsRepo   sms.Repo
	publisher sms.Publisher
}

func NewSMSUseCase(smsRepo sms.Repo, publisher sms.Publisher, db *gorm.DB) *UseCase {
	return &UseCase{
		smsRepo:   smsRepo.WithTx(db),
		publisher: publisher,
	}
}

func (u *UseCase) GetSMSByID(ctx context.Context, filter sms.Filter) (*sms.SMSMessage, error) {
	return u.smsRepo.GetByFilter(ctx, filter)
}

func (u *UseCase) ProcessSMS(ctx context.Context, smsMsg *sms.SMSMessage) error {
	err := u.smsRepo.Create(ctx, smsMsg)
	if err != nil {
		return err
	}

	debitEvent := sms.DebitUserBalance{
		UserID:    smsMsg.UserID,
		SMSID:     smsMsg.ID,
		Amount:    1,
		TimeStamp: time.Now(),
	}
	// TODO: use OutBox or update the sms failed.
	err = u.publisher.PublishEvent(ctx, debitEvent)
	if err != nil {
		return err
	}

	return nil
}
