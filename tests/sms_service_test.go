package tests

import (
	"context"
	"errors"
	"sms/internal/domain/sms"
	smsService "sms/internal/usecase/sms"
	"sms/pkg/logger"
	"testing"
	"time"

	"gorm.io/gorm"
)

type mockSMSRepo struct {
	messages    map[string]*sms.SMSMessage
	createError error
	getError    error
	updateError error
}

func newMockSMSRepo() *mockSMSRepo {
	return &mockSMSRepo{
		messages: make(map[string]*sms.SMSMessage),
	}
}

func (m *mockSMSRepo) GetByFilter(ctx context.Context, filter sms.Filter) (*sms.SMSMessage, error) {
	if m.getError != nil {
		return nil, m.getError
	}

	if filter.ID != nil {
		if msg, exists := m.messages[*filter.ID]; exists {
			return msg, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	for _, msg := range m.messages {
		if filter.Status != nil && msg.Status != *filter.Status {
			continue
		}
		if filter.UserID != nil && msg.UserID != *filter.UserID {
			continue
		}
		return msg, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (m *mockSMSRepo) Create(ctx context.Context, message *sms.SMSMessage) error {
	if m.createError != nil {
		return m.createError
	}
	m.messages[message.ID] = message
	return nil
}

func (m *mockSMSRepo) Update(ctx context.Context, ID string, message *sms.SMSMessage) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.messages[ID] = message
	return nil
}

func (m *mockSMSRepo) WithTx(tx *gorm.DB) sms.Repo {
	return m
}

type mockEventPublisher struct {
	publishedEvents []sms.DomainEvent
	publishError    error
}

func newMockEventPublisher() *mockEventPublisher {
	return &mockEventPublisher{
		publishedEvents: make([]sms.DomainEvent, 0),
	}
}

func (m *mockEventPublisher) PublishEvent(ctx context.Context, event sms.DomainEvent) error {
	if m.publishError != nil {
		return m.publishError
	}
	m.publishedEvents = append(m.publishedEvents, event)
	return nil
}

type mockSMSProvider struct {
	sendError    error
	providerName string
}

func newMockSMSProvider() *mockSMSProvider {
	return &mockSMSProvider{
		providerName: "mock-provider",
	}
}

func (m *mockSMSProvider) SendSMS(ctx context.Context, message *sms.SMSMessage) (string, error) {
	if m.sendError != nil {
		return m.providerName, m.sendError
	}
	return m.providerName, nil
}

func TestSMSService_CreateAndBillSMS_Success(t *testing.T) {
	repo := newMockSMSRepo()
	publisher := newMockEventPublisher()
	provider := newMockSMSProvider()
	log := logger.NewLogger("info")
	service := smsService.NewSMSService(repo, publisher, provider, &gorm.DB{}, log)

	message := &sms.SMSMessage{
		ID:       "test-sms-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}

	ctx := context.Background()

	err := service.CreateAndBillSMS(ctx, message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if _, exists := repo.messages[message.ID]; !exists {
		t.Error("Expected message to be created in repository")
	}

	if len(publisher.publishedEvents) != 1 {
		t.Errorf("Expected 1 published event, got %d", len(publisher.publishedEvents))
	}

	billingEvent, ok := publisher.publishedEvents[0].(sms.RequestSMSBilling)
	if !ok {
		t.Error("Expected published event to be RequestSMSBilling")
	}

	if billingEvent.UserID != message.UserID {
		t.Errorf("Expected billing event UserID to be %s, got %s", message.UserID, billingEvent.UserID)
	}
	if billingEvent.SMSID != message.ID {
		t.Errorf("Expected billing event SMSID to be %s, got %s", message.ID, billingEvent.SMSID)
	}
}

func TestSMSService_CreateAndBillSMS_RepoError(t *testing.T) {
	repo := newMockSMSRepo()
	repo.createError = errors.New("database error")
	publisher := newMockEventPublisher()
	provider := newMockSMSProvider()
	log := logger.NewLogger("info")
	service := smsService.NewSMSService(repo, publisher, provider, &gorm.DB{}, log)

	message := &sms.SMSMessage{
		ID:       "test-sms-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}

	ctx := context.Background()

	err := service.CreateAndBillSMS(ctx, message)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if len(publisher.publishedEvents) != 0 {
		t.Errorf("Expected no published events, got %d", len(publisher.publishedEvents))
	}
}

func TestSMSService_CreateAndBillSMS_PublishError(t *testing.T) {
	repo := newMockSMSRepo()
	publisher := newMockEventPublisher()
	publisher.publishError = errors.New("publish error")
	provider := newMockSMSProvider()
	log := logger.NewLogger("info")
	service := smsService.NewSMSService(repo, publisher, provider, &gorm.DB{}, log)

	message := &sms.SMSMessage{
		ID:       "test-sms-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}

	ctx := context.Background()

	err := service.CreateAndBillSMS(ctx, message)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, exists := repo.messages[message.ID]; !exists {
		t.Error("Expected message to be created in repository")
	}
}

func TestSMSService_ProcessDebitedSMS_Success(t *testing.T) {
	repo := newMockSMSRepo()
	publisher := newMockEventPublisher()
	provider := newMockSMSProvider()
	log := logger.NewLogger("info")
	service := smsService.NewSMSService(repo, publisher, provider, &gorm.DB{}, log)

	message := &sms.SMSMessage{
		ID:       "test-sms-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}
	repo.messages[message.ID] = message

	event := sms.SMSBillingCompleted{
		UserID:        "user-123",
		SMSID:         "test-sms-id",
		Amount:        1,
		TransactionID: "txn-123",
		TimeStamp:     time.Now(),
	}

	ctx := context.Background()

	err := service.ProcessDebitedSMS(ctx, event)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	updatedMessage := repo.messages[message.ID]
	if updatedMessage.Status != sms.SMSStatusDelivered {
		t.Errorf("Expected message status to be %s, got %s", sms.SMSStatusDelivered, updatedMessage.Status)
	}
	if updatedMessage.Provider != provider.providerName {
		t.Errorf("Expected provider to be %s, got %s", provider.providerName, updatedMessage.Provider)
	}
}

func TestSMSService_ProcessDebitedSMS_DeliveryFailure(t *testing.T) {
	repo := newMockSMSRepo()
	publisher := newMockEventPublisher()
	provider := newMockSMSProvider()
	provider.sendError = errors.New("network error")
	log := logger.NewLogger("info")
	service := smsService.NewSMSService(repo, publisher, provider, &gorm.DB{}, log)

	message := &sms.SMSMessage{
		ID:       "test-sms-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}
	repo.messages[message.ID] = message

	event := sms.SMSBillingCompleted{
		UserID:        "user-123",
		SMSID:         "test-sms-id",
		Amount:        1,
		TransactionID: "txn-123",
		TimeStamp:     time.Now(),
	}

	ctx := context.Background()

	err := service.ProcessDebitedSMS(ctx, event)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	updatedMessage := repo.messages[message.ID]
	if updatedMessage.Status != sms.SMSStatusFailed {
		t.Errorf("Expected message status to be %s, got %s", sms.SMSStatusFailed, updatedMessage.Status)
	}
	if updatedMessage.FailureCode != sms.MNOProviderFailed {
		t.Errorf("Expected failure code to be %s, got %s", sms.MNOProviderFailed, updatedMessage.FailureCode)
	}

	if len(publisher.publishedEvents) != 1 {
		t.Errorf("Expected 1 published event, got %d", len(publisher.publishedEvents))
	}

	refundEvent, ok := publisher.publishedEvents[0].(sms.RequestBillingRefund)
	if !ok {
		t.Error("Expected published event to be RequestBillingRefund")
	}
	if refundEvent.TransactionID != event.TransactionID {
		t.Errorf("Expected refund TransactionID to be %s, got %s", event.TransactionID, refundEvent.TransactionID)
	}
}

func TestSMSService_GetSMSByID_NotFound(t *testing.T) {
	repo := newMockSMSRepo()
	publisher := newMockEventPublisher()
	provider := newMockSMSProvider()
	log := logger.NewLogger("info")
	service := smsService.NewSMSService(repo, publisher, provider, &gorm.DB{}, log)

	nonExistentID := "non-existent-id"
	filter := sms.Filter{ID: &nonExistentID}
	ctx := context.Background()

	result, err := service.GetSMSByID(ctx, filter)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected no message to be returned")
	}
}
