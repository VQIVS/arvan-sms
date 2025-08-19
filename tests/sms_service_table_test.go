package tests

import (
	"context"
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"testing"
	"time"
)

func TestSMSService_SendSMS_TableDriven(t *testing.T) {
	tests := []struct {
		name                   string
		request                *presenter.SendSMSReq
		mockCreateSMSResponse  domain.SMSID
		mockCreateSMSError     error
		mockBalanceUpdateError error
		expectedResponse       *presenter.SendSMSResp
		expectedError          error
	}{
		{
			name: "successful SMS creation with valid data",
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+1234567890",
				Message:   "Hello World",
			},
			mockCreateSMSResponse:  domain.SMSID(123),
			mockCreateSMSError:     nil,
			mockBalanceUpdateError: nil,
			expectedResponse: &presenter.SendSMSResp{
				Status:  presenter.Pending,
				Message: "SMS created successfully",
			},
			expectedError: nil,
		},
		{
			name: "SMS creation fails",
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+1234567890",
				Message:   "Hello World",
			},
			mockCreateSMSResponse:  domain.SMSID(0),
			mockCreateSMSError:     &MockError{message: "database error"},
			mockBalanceUpdateError: nil,
			expectedResponse:       nil,
			expectedError:          &MockError{message: "database error"},
		},
		{
			name: "balance update fails",
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+1234567890",
				Message:   "Hello World",
			},
			mockCreateSMSResponse:  domain.SMSID(123),
			mockCreateSMSError:     nil,
			mockBalanceUpdateError: &MockError{message: "balance service error"},
			expectedResponse:       nil,
			expectedError:          &MockError{message: "balance service error"},
		},
		{
			name: "long message content",
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+1234567890",
				Message:   "This is a very long message that exceeds the typical SMS length limit to test how the service handles long messages and whether it processes them correctly without any issues.",
			},
			mockCreateSMSResponse:  domain.SMSID(456),
			mockCreateSMSError:     nil,
			mockBalanceUpdateError: nil,
			expectedResponse: &presenter.SendSMSResp{
				Status:  presenter.Pending,
				Message: "SMS created successfully",
			},
			expectedError: nil,
		},
		{
			name: "international phone number",
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+44123456789",
				Message:   "International SMS",
			},
			mockCreateSMSResponse:  domain.SMSID(789),
			mockCreateSMSError:     nil,
			mockBalanceUpdateError: nil,
			expectedResponse: &presenter.SendSMSResp{
				Status:  presenter.Pending,
				Message: "SMS created successfully",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{}
			smsService := service.NewSMSService(mockSvc)

			mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
				return tt.mockCreateSMSResponse, tt.mockCreateSMSError
			}

			mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
				return tt.mockBalanceUpdateError
			}

			resp, err := smsService.SendSMS(context.Background(), tt.request)
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.expectedResponse != nil {
				if resp == nil {
					t.Fatal("Expected response, got nil")
				}
				if resp.Status != tt.expectedResponse.Status {
					t.Errorf("Expected status %s, got %s", tt.expectedResponse.Status, resp.Status)
				}
				if resp.Message != tt.expectedResponse.Message {
					t.Errorf("Expected message %s, got %s", tt.expectedResponse.Message, resp.Message)
				}
			} else if resp != nil {
				t.Errorf("Expected nil response, got %v", resp)
			}
		})
	}
}

func TestSMSService_GetSMSMessage_TableDriven(t *testing.T) {
	tests := []struct {
		name             string
		smsID            uint
		mockSMSResponse  *domain.SMS
		mockError        error
		expectedResponse *presenter.SMSResp
		expectedError    error
	}{
		{
			name:  "successful SMS retrieval",
			smsID: 123,
			mockSMSResponse: &domain.SMS{
				ID:        domain.SMSID(123),
				Recipient: "+1234567890",
				Message:   "Hello World",
				Status:    string(domain.Delivered),
				CreatedAt: time.Now(),
			},
			mockError: nil,
			expectedResponse: &presenter.SMSResp{
				ID:        123,
				Recipient: "+1234567890",
				Message:   "Hello World",
				Status:    presenter.Status(domain.Delivered),
			},
			expectedError: nil,
		},
		{
			name:             "SMS not found",
			smsID:            999,
			mockSMSResponse:  nil,
			mockError:        &MockError{message: "SMS not found"},
			expectedResponse: nil,
			expectedError:    &MockError{message: "SMS not found"},
		},
		{
			name:  "SMS with pending status",
			smsID: 456,
			mockSMSResponse: &domain.SMS{
				ID:        domain.SMSID(456),
				Recipient: "+9876543210",
				Message:   "Pending message",
				Status:    string(domain.Pending),
				CreatedAt: time.Now(),
			},
			mockError: nil,
			expectedResponse: &presenter.SMSResp{
				ID:        456,
				Recipient: "+9876543210",
				Message:   "Pending message",
				Status:    presenter.Status(domain.Pending),
			},
			expectedError: nil,
		},
		{
			name:  "SMS with failed status",
			smsID: 789,
			mockSMSResponse: &domain.SMS{
				ID:        domain.SMSID(789),
				Recipient: "+1122334455",
				Message:   "Failed message",
				Status:    string(domain.Failed),
				CreatedAt: time.Now(),
			},
			mockError: nil,
			expectedResponse: &presenter.SMSResp{
				ID:        789,
				Recipient: "+1122334455",
				Message:   "Failed message",
				Status:    presenter.Status(domain.Failed),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{}
			smsService := service.NewSMSService(mockSvc)

			mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
				if filter.ID != domain.SMSID(tt.smsID) {
					t.Errorf("Expected SMS ID %d, got %d", tt.smsID, filter.ID)
				}
				return tt.mockSMSResponse, tt.mockError
			}

			resp, err := smsService.GetSMSMessage(context.Background(), tt.smsID)
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.expectedResponse != nil {
				if resp == nil {
					t.Fatal("Expected response, got nil")
				}
				if resp.ID != tt.expectedResponse.ID {
					t.Errorf("Expected ID %d, got %d", tt.expectedResponse.ID, resp.ID)
				}
				if resp.Recipient != tt.expectedResponse.Recipient {
					t.Errorf("Expected recipient %s, got %s", tt.expectedResponse.Recipient, resp.Recipient)
				}
				if resp.Message != tt.expectedResponse.Message {
					t.Errorf("Expected message %s, got %s", tt.expectedResponse.Message, resp.Message)
				}
				if resp.Status != tt.expectedResponse.Status {
					t.Errorf("Expected status %s, got %s", tt.expectedResponse.Status, resp.Status)
				}
			} else if resp != nil {
				t.Errorf("Expected nil response, got %v", resp)
			}
		})
	}
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}
