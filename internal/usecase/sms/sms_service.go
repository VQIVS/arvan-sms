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
	log       *logger.Logger
}

func NewSMSService(smsRepo sms.Repo, publisher sms.EventPublisher, provider sms.SMSProvider, db *gorm.DB, log *logger.Logger) *Service {
	return &Service{
		smsRepo:   smsRepo.WithTx(db),
		publisher: publisher,
		provider:  provider,
		log:       log,
	}
}

func (u *Service) GetSMSByID(ctx context.Context, filter sms.Filter) (*sms.SMSMessage, error) {
	return u.smsRepo.GetByFilter(ctx, filter)
}

func (u *Service) CreateAndBillSMS(ctx context.Context, smsMsg *sms.SMSMessage) error {
	u.log.Info(ctx, "creating SMS and requesting billing", "sms_id", smsMsg.ID, "user_id", smsMsg.UserID, "receiver", smsMsg.Receiver)

	err := u.smsRepo.Create(ctx, smsMsg)
	if err != nil {
		u.log.Error(ctx, "failed to create SMS in database", "error", err, "sms_id", smsMsg.ID)
		return err
	}
	u.log.Info(ctx, "SMS created successfully", "sms_id", smsMsg.ID)

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
		u.log.Error(ctx, "failed to publish billing request", "error", err, "sms_id", smsMsg.ID)
		return err
	}
	u.log.Info(ctx, "billing request published successfully", "sms_id", smsMsg.ID)

	return nil
}

func (u *Service) ProcessDebitedSMS(ctx context.Context, event sms.SMSBillingCompleted) error {
	u.log.Info(ctx, "processing debited SMS", "sms_id", event.SMSID, "transaction_id", event.TransactionID)

	smsMsg, err := u.smsRepo.GetByFilter(ctx, sms.Filter{ID: &event.SMSID})
	if err != nil {
		u.log.Error(ctx, "failed to retrieve SMS from database", "error", err, "sms_id", event.SMSID)
		return err
	}

	u.log.Info(ctx, "attempting SMS delivery", "sms_id", event.SMSID, "receiver", smsMsg.Receiver)
	provider, err := u.dispatchSMSDelivery(ctx, *smsMsg)

	if err != nil {
		u.log.Error(ctx, "SMS delivery failed", "error", err, "sms_id", event.SMSID, "provider", provider)
		smsMsg.MarkAsFailed(provider, sms.MNOProviderFailed)

		// refunding user
		refundMsg := sms.RequestBillingRefund{
			TransactionID: event.TransactionID,
			TimeStamp:     time.Now(),
		}

		u.log.Info(ctx, "publishing refund request due to delivery failure", "sms_id", event.SMSID, "transaction_id", event.TransactionID)
		err = u.publisher.PublishEvent(ctx, refundMsg)
		if err != nil {
			u.log.Error(ctx, "failed to publish refund request", "error", err, "sms_id", event.SMSID, "transaction_id", event.TransactionID)
			return err
		}
		u.log.Info(ctx, "refund request published successfully", "sms_id", event.SMSID, "transaction_id", event.TransactionID)
	} else {
		u.log.Info(ctx, "SMS delivered successfully", "sms_id", event.SMSID, "provider", provider, "receiver", smsMsg.Receiver)
		smsMsg.MarkAsSent(provider)
	}

	// updating sms object
	err = u.smsRepo.Update(ctx, smsMsg.ID, smsMsg)
	if err != nil {
		u.log.Error(ctx, "failed to update SMS status in database", "error", err, "sms_id", event.SMSID)
		return err
	}

	u.log.Info(ctx, "SMS processing completed", "sms_id", event.SMSID, "final_status", string(smsMsg.Status))
	return nil
}
func (u *Service) dispatchSMSDelivery(ctx context.Context, message sms.SMSMessage) (string, error) {
	return u.provider.SendSMS(ctx, &message)
}
