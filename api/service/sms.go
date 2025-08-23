package service

import (
	"context"
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"sms-dispatcher/internal/sms/port"
)

type SMSService struct {
	svc port.Service
}

func NewSMSService(svc port.Service) *SMSService {
	return &SMSService{
		svc: svc,
	}
}

func (s *SMSService) SendSMS(ctx context.Context, req *presenter.SendSMSReq) (*presenter.SendSMSResp, error) {
	id, err := s.svc.CreateSMS(ctx, req.Recipient, req.Message)
	if err != nil {
		return nil, err
	}
	err = s.svc.UserBalanceUpdate(
		ctx,
		event.UserBalanceEvent{
			Domain: event.SMS,
			UserID: req.UserID,
			SMSID:  uint(id),
			//FIXME: Define amount in some where better like config
			Amount: 1,
			Type:   event.SMSCreditEvent,
		},
	)

	if err != nil {
		return nil, err
	}

	return &presenter.SendSMSResp{
		ID:      uint(id),
		Status:  presenter.Pending,
		Message: "SMS created successfully",
	}, nil
}

func (s *SMSService) GetSMSMessage(ctx context.Context, ID uint) (*presenter.SMSResp, error) {
	sms, err := s.svc.GetSMSByFilter(ctx, &domain.SMSFilter{
		ID: domain.SMSID(ID),
	})
	if err != nil {
		return nil, err
	}

	return &presenter.SMSResp{
		ID:        uint(sms.ID),
		Recipient: sms.Recipient,
		Message:   sms.Message,
		Status:    presenter.Status(sms.Status),
	}, nil
}
