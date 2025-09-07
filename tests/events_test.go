package tests

import (
	"context"
	"sms/internal/domain/sms"
	"testing"
	"time"
)

func TestRequestSMSBilling_Event(t *testing.T) {
	userID := "user-123"
	smsID := "sms-456"
	amount := int64(100)
	timestamp := time.Now()

	event := sms.RequestSMSBilling{
		UserID:    userID,
		SMSID:     smsID,
		Amount:    amount,
		TimeStamp: timestamp,
	}

	if event.EventType() != sms.EventTypeBillingRequested {
		t.Errorf("Expected event type to be %s, got %s", sms.EventTypeBillingRequested, event.EventType())
	}
	if event.AggregateID() != smsID {
		t.Errorf("Expected aggregate ID to be %s, got %s", smsID, event.AggregateID())
	}
	if !event.Timestamp().Equal(timestamp) {
		t.Errorf("Expected timestamp to be %v, got %v", timestamp, event.Timestamp())
	}
	if event.UserID != userID {
		t.Errorf("Expected UserID to be %s, got %s", userID, event.UserID)
	}
	if event.SMSID != smsID {
		t.Errorf("Expected SMSID to be %s, got %s", smsID, event.SMSID)
	}
	if event.Amount != amount {
		t.Errorf("Expected Amount to be %d, got %d", amount, event.Amount)
	}
}

func TestSMSBillingCompleted_Event(t *testing.T) {
	userID := "user-123"
	smsID := "sms-456"
	amount := int64(100)
	transactionID := "txn-789"
	timestamp := time.Now()

	event := sms.SMSBillingCompleted{
		UserID:        userID,
		SMSID:         smsID,
		Amount:        amount,
		TransactionID: transactionID,
		TimeStamp:     timestamp,
	}

	if event.EventType() != sms.EventTypeBillingCompleted {
		t.Errorf("Expected event type to be %s, got %s", sms.EventTypeBillingCompleted, event.EventType())
	}
	if event.AggregateID() != smsID {
		t.Errorf("Expected aggregate ID to be %s, got %s", smsID, event.AggregateID())
	}
	if !event.Timestamp().Equal(timestamp) {
		t.Errorf("Expected timestamp to be %v, got %v", timestamp, event.Timestamp())
	}
	if event.UserID != userID {
		t.Errorf("Expected UserID to be %s, got %s", userID, event.UserID)
	}
	if event.SMSID != smsID {
		t.Errorf("Expected SMSID to be %s, got %s", smsID, event.SMSID)
	}
	if event.Amount != amount {
		t.Errorf("Expected Amount to be %d, got %d", amount, event.Amount)
	}
	if event.TransactionID != transactionID {
		t.Errorf("Expected TransactionID to be %s, got %s", transactionID, event.TransactionID)
	}
}

func TestRequestBillingRefund_Event(t *testing.T) {
	transactionID := "txn-789"
	timestamp := time.Now()

	event := sms.RequestBillingRefund{
		TransactionID: transactionID,
		TimeStamp:     timestamp,
	}

	if event.EventType() != sms.EventTypeBillingRefunded {
		t.Errorf("Expected event type to be %s, got %s", sms.EventTypeBillingRefunded, event.EventType())
	}
	if event.AggregateID() != transactionID {
		t.Errorf("Expected aggregate ID to be %s, got %s", transactionID, event.AggregateID())
	}
	if !event.Timestamp().Equal(timestamp) {
		t.Errorf("Expected timestamp to be %v, got %v", timestamp, event.Timestamp())
	}
	if event.TransactionID != transactionID {
		t.Errorf("Expected TransactionID to be %s, got %s", transactionID, event.TransactionID)
	}
}

func TestEventType_Constants(t *testing.T) {
	tests := []struct {
		name      string
		eventType sms.EventType
		expected  string
	}{
		{"Billing Requested", sms.EventTypeBillingRequested, "BillingRequested"},
		{"Billing Completed", sms.EventTypeBillingCompleted, "BillingCompleted"},
		{"Billing Refunded", sms.EventTypeBillingRefunded, "BillingRefunded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.eventType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.eventType))
			}
		})
	}
}

func TestSMSProviderFunc(t *testing.T) {
	providerFunc := sms.SMSProviderFunc(func(ctx context.Context, message *sms.SMSMessage) (string, error) {
		return "test-provider", nil
	})

	var provider sms.SMSProvider = providerFunc

	message := &sms.SMSMessage{
		ID:       "test-id",
		Content:  "test",
		Receiver: "+1234567890",
	}

	result, err := provider.SendSMS(context.Background(), message)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "test-provider" {
		t.Errorf("Expected provider name to be 'test-provider', got %s", result)
	}
}
