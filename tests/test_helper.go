package tests

import (
	"context"
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"time"
)

type TestHelper struct{}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (h *TestHelper) CreateDefaultSMSService() (*service.SMSService, *MockService) {
	mockSvc := &MockService{}

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		return &domain.SMS{
			ID:        filter.ID,
			Recipient: "+1234567890",
			Message:   "Test message",
			Status:    string(domain.Delivered),
			CreatedAt: time.Now(),
		}, nil
	}

	return service.NewSMSService(mockSvc), mockSvc
}

func (h *TestHelper) CreateSampleSendSMSRequest() *presenter.SendSMSReq {
	return &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Sample test message",
	}
}

func (h *TestHelper) CreateSampleSMS(id uint) *domain.SMS {
	return &domain.SMS{
		ID:        domain.SMSID(id),
		Recipient: "+1234567890",
		Message:   "Sample SMS message",
		Status:    string(domain.Delivered),
		CreatedAt: time.Now(),
	}
}

func (h *TestHelper) CreateUserBalanceEvent(userID uint, smsID uint) event.UserBalanceEvent {
	return event.UserBalanceEvent{
		Domain: event.SMS,
		UserID: userID,
		SMSID:  smsID,
		Amount: 1,
		Type:   event.SMSCreditEvent,
	}
}

func (h *TestHelper) AssertSendSMSResponse(t interface{}, resp *presenter.SendSMSResp, expectedStatus presenter.Status, expectedMessage string) {
	type TestingT interface {
		Errorf(format string, args ...interface{})
		Fatal(args ...interface{})
	}

	test := t.(TestingT)

	if resp == nil {
		test.Fatal("Expected response, got nil")
		return
	}

	if resp.Status != expectedStatus {
		test.Errorf("Expected status %s, got %s", expectedStatus, resp.Status)
	}

	if resp.Message != expectedMessage {
		test.Errorf("Expected message %s, got %s", expectedMessage, resp.Message)
	}
}

func (h *TestHelper) AssertSMSResponse(t interface{}, resp *presenter.SMSResp, expectedID uint, expectedRecipient string, expectedMessage string, expectedStatus presenter.Status) {
	type TestingT interface {
		Errorf(format string, args ...interface{})
		Fatal(args ...interface{})
	}

	test := t.(TestingT)

	if resp == nil {
		test.Fatal("Expected response, got nil")
		return
	}

	if resp.ID != expectedID {
		test.Errorf("Expected ID %d, got %d", expectedID, resp.ID)
	}

	if resp.Recipient != expectedRecipient {
		test.Errorf("Expected recipient %s, got %s", expectedRecipient, resp.Recipient)
	}

	if resp.Message != expectedMessage {
		test.Errorf("Expected message %s, got %s", expectedMessage, resp.Message)
	}

	if resp.Status != expectedStatus {
		test.Errorf("Expected status %s, got %s", expectedStatus, resp.Status)
	}
}

func (h *TestHelper) GenerateTestPhoneNumbers(count int) []string {
	numbers := make([]string, count)
	for i := 0; i < count; i++ {
		numbers[i] = "+123456789" + string(rune(48+i%10))
	}
	return numbers
}

func (h *TestHelper) GenerateTestMessages(count int) []string {
	messages := make([]string, count)
	for i := 0; i < count; i++ {
		messages[i] = "Test message " + string(rune(48+i%10))
	}
	return messages
}

const (
	DefaultTestUserID    = uint(1)
	DefaultTestRecipient = "+1234567890"
	DefaultTestMessage   = "Test message"
	DefaultTestSMSID     = uint(123)
)
