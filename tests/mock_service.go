package tests

import (
	"context"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
)

// MockService is a mock implementation of the port.Service interface
type MockService struct {
	CreateSMSFunc         func(ctx context.Context, recipient string, message string) (domain.SMSID, error)
	GetSMSByFilterFunc    func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error)
	UserBalanceUpdateFunc func(ctx context.Context, user event.UserBalanceEvent) error
	UpdateSMSStatusFunc   func(ctx context.Context, body []byte) error
}

func (m *MockService) CreateSMS(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
	if m.CreateSMSFunc != nil {
		return m.CreateSMSFunc(ctx, recipient, message)
	}
	return 0, nil
}

func (m *MockService) GetSMSByFilter(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
	if m.GetSMSByFilterFunc != nil {
		return m.GetSMSByFilterFunc(ctx, filter)
	}
	return nil, nil
}

func (m *MockService) UserBalanceUpdate(ctx context.Context, user event.UserBalanceEvent) error {
	if m.UserBalanceUpdateFunc != nil {
		return m.UserBalanceUpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockService) UpdateSMSStatus(ctx context.Context, body []byte) error {
	if m.UpdateSMSStatusFunc != nil {
		return m.UpdateSMSStatusFunc(ctx, body)
	}
	return nil
}
