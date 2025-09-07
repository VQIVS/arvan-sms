package tests

import (
	"testing"
	"time"

	"sms/internal/api/dto"
)

func TestSendSMSRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.SendSMSRequest
		isValid bool
	}{
		{
			name: "Valid SMS request",
			request: dto.SendSMSRequest{
				Content:  "Hello, this is a test message",
				Receiver: "+1234567890",
				UserID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			isValid: true,
		},
		{
			name: "Empty content",
			request: dto.SendSMSRequest{
				Content:  "",
				Receiver: "+1234567890",
				UserID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			isValid: false,
		},
		{
			name: "Content too long (over 160 chars)",
			request: dto.SendSMSRequest{
				Content:  "This is a very long message that exceeds the 160 character limit for SMS messages. It should fail validation because SMS messages have a strict character limit that must be enforced properly in all cases.",
				Receiver: "+1234567890",
				UserID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			isValid: false,
		},
		{
			name: "Empty receiver",
			request: dto.SendSMSRequest{
				Content:  "Test message",
				Receiver: "",
				UserID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			isValid: false,
		},
		{
			name: "Empty user ID",
			request: dto.SendSMSRequest{
				Content:  "Test message",
				Receiver: "+1234567890",
				UserID:   "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := true

			if tt.request.Content == "" {
				isValid = false
			}
			if len(tt.request.Content) > 160 {
				isValid = false
			}
			if tt.request.Receiver == "" {
				isValid = false
			}
			if tt.request.UserID == "" {
				isValid = false
			}

			if isValid != tt.isValid {
				t.Errorf("Expected validation result %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestSendSMSResponse_Creation(t *testing.T) {
	id := "test-sms-id"
	status := "pending"
	createdAt := time.Now()
	message := "SMS queued for delivery"

	response := dto.SendSMSResponse{
		ID:        id,
		Status:    status,
		CreatedAt: createdAt,
		Message:   message,
	}

	if response.ID != id {
		t.Errorf("Expected ID to be %s, got %s", id, response.ID)
	}
	if response.Status != status {
		t.Errorf("Expected status to be %s, got %s", status, response.Status)
	}
	if !response.CreatedAt.Equal(createdAt) {
		t.Errorf("Expected CreatedAt to be %v, got %v", createdAt, response.CreatedAt)
	}
	if response.Message != message {
		t.Errorf("Expected message to be %s, got %s", message, response.Message)
	}
}

func TestGetSMSResponse_Creation(t *testing.T) {
	id := "test-sms-id"
	userID := "user-123"
	content := "Test SMS content"
	receiver := "+1234567890"
	provider := "test-provider"
	status := "delivered"
	deliveredAt := time.Now()
	failureCode := ""
	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := time.Now()

	response := dto.GetSMSResponse{
		ID:          id,
		UserID:      userID,
		Content:     content,
		Receiver:    receiver,
		Provider:    provider,
		Status:      status,
		DeliveredAt: &deliveredAt,
		FailureCode: failureCode,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	if response.ID != id {
		t.Errorf("Expected ID to be %s, got %s", id, response.ID)
	}
	if response.UserID != userID {
		t.Errorf("Expected UserID to be %s, got %s", userID, response.UserID)
	}
	if response.Content != content {
		t.Errorf("Expected Content to be %s, got %s", content, response.Content)
	}
	if response.Receiver != receiver {
		t.Errorf("Expected Receiver to be %s, got %s", receiver, response.Receiver)
	}
	if response.Provider != provider {
		t.Errorf("Expected Provider to be %s, got %s", provider, response.Provider)
	}
	if response.Status != status {
		t.Errorf("Expected Status to be %s, got %s", status, response.Status)
	}
	if response.DeliveredAt == nil || !response.DeliveredAt.Equal(deliveredAt) {
		t.Errorf("Expected DeliveredAt to be %v, got %v", deliveredAt, response.DeliveredAt)
	}
	if response.FailureCode != failureCode {
		t.Errorf("Expected FailureCode to be %s, got %s", failureCode, response.FailureCode)
	}
}

func TestGetSMSResponse_FailedSMS(t *testing.T) {
	response := dto.GetSMSResponse{
		ID:          "failed-sms-id",
		UserID:      "user-123",
		Content:     "Test SMS content",
		Receiver:    "+1234567890",
		Provider:    "test-provider",
		Status:      "failed",
		DeliveredAt: nil,
		FailureCode: "NETWORK_ERROR",
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	if response.Status != "failed" {
		t.Errorf("Expected status to be failed, got %s", response.Status)
	}
	if response.DeliveredAt != nil {
		t.Error("Expected DeliveredAt to be nil for failed SMS")
	}
	if response.FailureCode != "NETWORK_ERROR" {
		t.Errorf("Expected FailureCode to be NETWORK_ERROR, got %s", response.FailureCode)
	}
}

func TestErrorResponse_Creation(t *testing.T) {
	errorMsg := "Validation failed"
	message := "Content cannot be empty"
	code := 400

	response := dto.ErrorResponse{
		Error:   errorMsg,
		Message: message,
		Code:    code,
	}

	if response.Error != errorMsg {
		t.Errorf("Expected Error to be %s, got %s", errorMsg, response.Error)
	}
	if response.Message != message {
		t.Errorf("Expected Message to be %s, got %s", message, response.Message)
	}
	if response.Code != code {
		t.Errorf("Expected Code to be %d, got %d", code, response.Code)
	}
}
