package tests

import (
	"sms/internal/domain/sms"
	"testing"
	"time"
)

func TestSMSMessage_MarkAsSent(t *testing.T) {
	message := &sms.SMSMessage{
		ID:       "test-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}
	provider := "test-provider"
	beforeTime := time.Now()

	message.MarkAsSent(provider)
	if message.Status != sms.SMSStatusDelivered {
		t.Errorf("Expected status to be %s, got %s", sms.SMSStatusDelivered, message.Status)
	}
	if message.Provider != provider {
		t.Errorf("Expected provider to be %s, got %s", provider, message.Provider)
	}
	if message.DeliveredAt.Before(beforeTime) {
		t.Error("DeliveredAt should be set to current time")
	}
	if message.UpdatedAt.Before(beforeTime) {
		t.Error("UpdatedAt should be set to current time")
	}
}

func TestSMSMessage_MarkAsFailed(t *testing.T) {
	message := &sms.SMSMessage{
		ID:       "test-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}
	provider := "test-provider"
	failureCode := "NETWORK_ERROR"
	beforeTime := time.Now()

	message.MarkAsFailed(provider, failureCode)
	if message.Status != sms.SMSStatusFailed {
		t.Errorf("Expected status to be %s, got %s", sms.SMSStatusFailed, message.Status)
	}
	if message.Provider != provider {
		t.Errorf("Expected provider to be %s, got %s", provider, message.Provider)
	}
	if message.FailureCode != failureCode {
		t.Errorf("Expected failure code to be %s, got %s", failureCode, message.FailureCode)
	}
	if message.UpdatedAt.Before(beforeTime) {
		t.Error("UpdatedAt should be set to current time")
	}
}

func TestSMSStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   sms.SMSStatus
		expected string
	}{
		{"Pending status", sms.SMSStatusPending, "pending"},
		{"Delivered status", sms.SMSStatusDelivered, "delivered"},
		{"Failed status", sms.SMSStatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.status))
			}
		})
	}
}

func TestFilter(t *testing.T) {
	id := "test-id"
	status := sms.SMSStatusDelivered
	userID := "user-123"

	tests := []struct {
		name      string
		filter    sms.Filter
		hasID     bool
		hasStatus bool
		hasUserID bool
	}{
		{
			name:   "Empty filter",
			filter: sms.Filter{},
		},
		{
			name:   "Filter with ID only",
			filter: sms.Filter{ID: &id},
			hasID:  true,
		},
		{
			name:      "Filter with status only",
			filter:    sms.Filter{Status: &status},
			hasStatus: true,
		},
		{
			name:      "Filter with user ID only",
			filter:    sms.Filter{UserID: &userID},
			hasUserID: true,
		},
		{
			name:      "Filter with all fields",
			filter:    sms.Filter{ID: &id, Status: &status, UserID: &userID},
			hasID:     true,
			hasStatus: true,
			hasUserID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasID && (tt.filter.ID == nil || *tt.filter.ID != id) {
				t.Error("Expected filter to have ID")
			}
			if !tt.hasID && tt.filter.ID != nil {
				t.Error("Expected filter to not have ID")
			}

			if tt.hasStatus && (tt.filter.Status == nil || *tt.filter.Status != status) {
				t.Error("Expected filter to have status")
			}
			if !tt.hasStatus && tt.filter.Status != nil {
				t.Error("Expected filter to not have status")
			}

			if tt.hasUserID && (tt.filter.UserID == nil || *tt.filter.UserID != userID) {
				t.Error("Expected filter to have user ID")
			}
			if !tt.hasUserID && tt.filter.UserID != nil {
				t.Error("Expected filter to not have user ID")
			}
		})
	}
}
