package sms

import (
	"context"
	"sms/internal/domain/sms"
	"sms/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	smsRepo   sms.Repo
	publisher sms.EventPublisher
	provider  sms.SMSProvider
}

func NewSMSService(smsRepo sms.Repo, publisher sms.EventPublisher, provider sms.SMSProvider, db *gorm.DB) *Service {
	return &Service{
		smsRepo:   smsRepo.WithTx(db),
		publisher: publisher,
		provider:  provider,
	}
}

func (u *Service) GetSMSByID(ctx context.Context, filter sms.Filter) (*sms.SMSMessage, error) {
	return u.smsRepo.GetByFilter(ctx, filter)
}

func (u *Service) CreateAndBillSMS(ctx context.Context, smsMsg *sms.SMSMessage) error {
	err := u.smsRepo.Create(ctx, smsMsg)
	logger.NewLogger("").ErrorWithoutContext("after create sms " + smsMsg.ID)
	if err != nil {
		return err
	}

	debitEvent := sms.RequestSMSBilling{
		UserID: smsMsg.UserID,
		SMSID:  smsMsg.ID,
		//TODO: do not hardcode amount
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

func (u *Service) ProcessDebitedSMS(ctx context.Context, event sms.SMSBillingCompleted) error {
	smsMsg, err := u.smsRepo.GetByFilter(ctx, sms.Filter{ID: &event.SMSID})
	if err != nil {
		return err
	}
	provider, err := u.dispatchSMSDelivery(ctx, *smsMsg)
	if err != nil {
		smsMsg.MarkAsFailed(provider, sms.MNOProviderFailed)
		// refunding user
		refundMsg := sms.RequestBillingRefund{
			TransactionID: event.TransactionID,
			TimeStamp:     time.Now(),
		}
		err = u.publisher.PublishEvent(ctx, refundMsg)
		if err != nil {
			return err
		}
	}
	smsMsg.MarkAsSent(provider)

	// updating sent sms object
	err = u.smsRepo.Update(ctx, smsMsg.ID, smsMsg)
	if err != nil {
		return err
	}
	return nil
}
func (u *Service) dispatchSMSDelivery(ctx context.Context, message sms.SMSMessage) (string, error) {
	return u.provider.SendSMS(ctx, &message)
}
