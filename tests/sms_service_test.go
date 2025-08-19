package tests

import (
	"context"
	"errors"
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"testing"
	"time"
)

func TestSMSService_SendSMS_Success(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	expectedSMSID := domain.SMSID(123)
	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Test message",
	}

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		if recipient != req.Recipient {
			t.Errorf("Expected recipient %s, got %s", req.Recipient, recipient)
		}
		if message != req.Message {
			t.Errorf("Expected message %s, got %s", req.Message, message)
		}
		return expectedSMSID, nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		expectedEvent := event.UserBalanceEvent{
			Domain: event.SMS,
			UserID: req.UserID,
			SMSID:  uint(expectedSMSID),
			Amount: 1,
			Type:   event.SMSCreditEvent,
		}

		if user.Domain != expectedEvent.Domain {
			t.Errorf("Expected domain %s, got %s", expectedEvent.Domain, user.Domain)
		}
		if user.UserID != expectedEvent.UserID {
			t.Errorf("Expected UserID %d, got %d", expectedEvent.UserID, user.UserID)
		}
		if user.SMSID != expectedEvent.SMSID {
			t.Errorf("Expected SMSID %d, got %d", expectedEvent.SMSID, user.SMSID)
		}
		if user.Amount != expectedEvent.Amount {
			t.Errorf("Expected Amount %f, got %f", expectedEvent.Amount, user.Amount)
		}
		if user.Type != expectedEvent.Type {
			t.Errorf("Expected Type %s, got %s", expectedEvent.Type, user.Type)
		}

		return nil
	}

	resp, err := smsService.SendSMS(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Status != presenter.Pending {
		t.Errorf("Expected status %s, got %s", presenter.Pending, resp.Status)
	}

	if resp.Message != "SMS created successfully" {
		t.Errorf("Expected message 'SMS created successfully', got %s", resp.Message)
	}
}

func TestSMSService_SendSMS_CreateSMSError(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Test message",
	}

	expectedError := errors.New("failed to create SMS")
	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return 0, expectedError
	}

	resp, err := smsService.SendSMS(context.Background(), req)

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}
}

func TestSMSService_SendSMS_UserBalanceUpdateError(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Test message",
	}

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	expectedError := errors.New("failed to update user balance")
	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return expectedError
	}

	resp, err := smsService.SendSMS(context.Background(), req)

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}
}

func TestSMSService_GetSMSMessage_Success(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	smsID := uint(123)
	expectedSMS := &domain.SMS{
		ID:        domain.SMSID(smsID),
		Recipient: "+1234567890",
		Message:   "Test message",
		Status:    string(domain.Delivered),
		CreatedAt: time.Now(),
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		if filter.ID != domain.SMSID(smsID) {
			t.Errorf("Expected SMS ID %d, got %d", smsID, filter.ID)
		}
		return expectedSMS, nil
	}

	resp, err := smsService.GetSMSMessage(context.Background(), smsID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.ID != uint(expectedSMS.ID) {
		t.Errorf("Expected ID %d, got %d", expectedSMS.ID, resp.ID)
	}

	if resp.Recipient != expectedSMS.Recipient {
		t.Errorf("Expected recipient %s, got %s", expectedSMS.Recipient, resp.Recipient)
	}

	if resp.Message != expectedSMS.Message {
		t.Errorf("Expected message %s, got %s", expectedSMS.Message, resp.Message)
	}

	if resp.Status != presenter.Status(expectedSMS.Status) {
		t.Errorf("Expected status %s, got %s", expectedSMS.Status, resp.Status)
	}
}

func TestSMSService_GetSMSMessage_Error(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	smsID := uint(123)
	expectedError := errors.New("SMS not found")

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		return nil, expectedError
	}

	resp, err := smsService.GetSMSMessage(context.Background(), smsID)

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}
}

func TestNewSMSService(t *testing.T) {
	mockSvc := &MockService{}

	smsService := service.NewSMSService(mockSvc)

	if smsService == nil {
		t.Fatal("Expected SMS service instance, got nil")
	}
}

func TestSMSService_SendSMS_EmptyRecipient(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "",
		Message:   "Test message",
	}

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	resp, err := smsService.SendSMS(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.Status != presenter.Pending {
		t.Errorf("Expected status %s, got %s", presenter.Pending, resp.Status)
	}
}

func TestSMSService_SendSMS_EmptyMessage(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "",
	}

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	resp, err := smsService.SendSMS(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.Status != presenter.Pending {
		t.Errorf("Expected status %s, got %s", presenter.Pending, resp.Status)
	}
}

func TestSMSService_GetSMSMessage_ZeroID(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		if filter.ID != domain.SMSID(0) {
			t.Errorf("Expected SMS ID 0, got %d", filter.ID)
		}
		return &domain.SMS{
			ID:        domain.SMSID(0),
			Recipient: "+1234567890",
			Message:   "Test message",
			Status:    string(domain.Pending),
		}, nil
	}

	resp, err := smsService.GetSMSMessage(context.Background(), 0)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.ID != 0 {
		t.Errorf("Expected ID 0, got %d", resp.ID)
	}
}
