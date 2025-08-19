package tests

import (
	"context"
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"sync"
	"testing"
	"time"
)

func TestSMSService_ConcurrentRequests(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var createSMSCalls int32
	var balanceUpdateCalls int32
	var mu sync.Mutex

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		mu.Lock()
		createSMSCalls++
		smsID := domain.SMSID(createSMSCalls)
		mu.Unlock()

		time.Sleep(1 * time.Millisecond)
		return smsID, nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		mu.Lock()
		balanceUpdateCalls++
		mu.Unlock()

		time.Sleep(1 * time.Millisecond)
		return nil
	}

	const numConcurrent = 10
	var wg sync.WaitGroup
	results := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			req := &presenter.SendSMSReq{
				UserID:    uint(id + 1),
				Recipient: "+123456789" + string(rune(48+id)),
				Message:   "Concurrent test message",
			}

			_, err := smsService.SendSMS(context.Background(), req)
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	errorCount := 0
	for err := range results {
		if err != nil {
			t.Errorf("Concurrent request failed: %v", err)
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected no errors, but got %d errors", errorCount)
	}

	if createSMSCalls != numConcurrent {
		t.Errorf("Expected %d CreateSMS calls, got %d", numConcurrent, createSMSCalls)
	}

	if balanceUpdateCalls != numConcurrent {
		t.Errorf("Expected %d UserBalanceUpdate calls, got %d", numConcurrent, balanceUpdateCalls)
	}
}

func TestSMSService_ContextCancellation(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return domain.SMSID(123), nil
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Test message",
	}

	resp, err := smsService.SendSMS(ctx, req)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}
}

func TestSMSService_ContextTimeout(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		select {
		case <-time.After(100 * time.Millisecond):
			return domain.SMSID(123), nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Test message",
	}

	resp, err := smsService.SendSMS(ctx, req)

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}
}

func TestSMSService_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockService)
		request     *presenter.SendSMSReq
		expectError bool
	}{
		{
			name: "zero user ID",
			setupMock: func(m *MockService) {
				m.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
					return domain.SMSID(1), nil
				}
				m.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
					if user.UserID != 0 {
						t.Errorf("Expected UserID 0, got %d", user.UserID)
					}
					return nil
				}
			},
			request: &presenter.SendSMSReq{
				UserID:    0,
				Recipient: "+1234567890",
				Message:   "Test",
			},
			expectError: false,
		},
		{
			name: "very long recipient number",
			setupMock: func(m *MockService) {
				m.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
					return domain.SMSID(1), nil
				}
				m.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
					return nil
				}
			},
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+123456789012345678901234567890",
				Message:   "Test",
			},
			expectError: false,
		},
		{
			name: "message with special characters",
			setupMock: func(m *MockService) {
				m.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
					expectedMessage := "Test with émojis 🎉 and spëcial châractérs @#$%^&*()"
					if message != expectedMessage {
						t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
					}
					return domain.SMSID(1), nil
				}
				m.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
					return nil
				}
			},
			request: &presenter.SendSMSReq{
				UserID:    1,
				Recipient: "+1234567890",
				Message:   "Test with émojis 🎉 and spëcial châractérs @#$%^&*()",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{}
			tt.setupMock(mockSvc)
			smsService := service.NewSMSService(mockSvc)

			resp, err := smsService.SendSMS(context.Background(), tt.request)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError && resp == nil {
				t.Error("Expected response but got nil")
			}
		})
	}
}

func TestSMSService_GetSMSMessage_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		smsID       uint
		setupMock   func(*MockService)
		expectError bool
	}{
		{
			name:  "maximum uint ID",
			smsID: ^uint(0),
			setupMock: func(m *MockService) {
				m.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
					if filter.ID != domain.SMSID(^uint(0)) {
						t.Errorf("Expected SMS ID %d, got %d", ^uint(0), filter.ID)
					}
					return &domain.SMS{
						ID:        domain.SMSID(^uint(0)),
						Recipient: "+1234567890",
						Message:   "Test",
						Status:    string(domain.Delivered),
					}, nil
				}
			},
			expectError: false,
		},
		{
			name:  "SMS with empty status",
			smsID: 123,
			setupMock: func(m *MockService) {
				m.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
					return &domain.SMS{
						ID:        domain.SMSID(123),
						Recipient: "+1234567890",
						Message:   "Test",
						Status:    "", // Empty status
					}, nil
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{}
			tt.setupMock(mockSvc)
			smsService := service.NewSMSService(mockSvc)

			resp, err := smsService.GetSMSMessage(context.Background(), tt.smsID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError && resp == nil {
				t.Error("Expected response but got nil")
			}
		})
	}
}
