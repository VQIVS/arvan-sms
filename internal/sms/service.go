package sms

import (
	"context"
	"encoding/json"
	"errors"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"sms-dispatcher/internal/sms/port"
	mno "sms-dispatcher/pkg/adapters/mno_mock"
	"sms-dispatcher/pkg/adapters/rabbit"
	"sms-dispatcher/pkg/constants"
)

type service struct {
	repo   port.Repo
	rabbit *rabbit.Rabbit
}

func NewService(repo port.Repo, r *rabbit.Rabbit) port.Service {
	return &service{
		repo:   repo,
		rabbit: r,
	}
}

func (s *service) CreateSMS(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
	sms := domain.SMS{
		Recipient: recipient,
		Message:   message,
		Status:    string(domain.Pending),
	}
	smsID, err := s.repo.Create(ctx, sms)
	if err != nil {
		return 0, err
	}
	return smsID, nil
}
func (s *service) GetSMSByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
	sms, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return sms, nil
}
func (s *service) UserBalanceUpdate(ctx context.Context, user event.UserBalanceEvent) error {
	if s.rabbit == nil {
		return nil
	}
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	s.rabbit.Logger.Info("publishing user balance update", "userID", user.UserID, "amount", user.Amount)
	return s.rabbit.Publish(body, constants.KeyBalanceUpdate)
}

func (s *service) UpdateSMSStatus(ctx context.Context, body []byte) error {
	var sms event.SMSUpdateEvent
	if err := json.Unmarshal(body, &sms); err != nil {
		return err
	}
	smsDomain, err := s.repo.GetByFilter(ctx,
		&domain.SMSFilter{
			ID: domain.SMSID(sms.SMSID),
		},
	)
	if err != nil {
		return err
	}
	switch sms.Status {
	case event.StatusFailed:
		smsDomain.Status = string(sms.Status)
		return s.repo.Update(ctx, *smsDomain)
	case event.StatusSuccess:
		mnoStatus, err := mno.SendSMSViaMNO()
		if err != nil {
			return err
		}
		if mnoStatus != mno.SuccessCode {
			return errors.New("failed to send SMS via MNO, retrying")
		}
		smsDomain.Status = string(sms.Status)
	}

	return s.repo.Update(ctx, *smsDomain)

}
