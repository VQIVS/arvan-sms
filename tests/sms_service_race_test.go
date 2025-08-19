package tests

import (
	"context"
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/internal/sms/event"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSMSService_SendSMS_RaceCondition(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var smsIDCounter int64
	var balanceUpdateCounter int64

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		newID := atomic.AddInt64(&smsIDCounter, 1)
		time.Sleep(1 * time.Millisecond)
		return domain.SMSID(newID), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		atomic.AddInt64(&balanceUpdateCounter, 1)
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	const numGoroutines = 100
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			req := &presenter.SendSMSReq{
				UserID:    uint(id),
				Recipient: "+1234567890",
				Message:   "Race test message",
			}

			_, err := smsService.SendSMS(context.Background(), req)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Unexpected error in race test: %v", err)
	}

	if atomic.LoadInt64(&smsIDCounter) != numGoroutines {
		t.Errorf("Expected %d SMS creations, got %d", numGoroutines, smsIDCounter)
	}

	if atomic.LoadInt64(&balanceUpdateCounter) != numGoroutines {
		t.Errorf("Expected %d balance updates, got %d", numGoroutines, balanceUpdateCounter)
	}
}

func TestSMSService_GetSMSMessage_RaceCondition(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var accessCounter int64

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		atomic.AddInt64(&accessCounter, 1)
		time.Sleep(1 * time.Millisecond)

		return &domain.SMS{
			ID:        filter.ID,
			Recipient: "+1234567890",
			Message:   "Race test message",
			Status:    string(domain.Delivered),
			CreatedAt: time.Now(),
		}, nil
	}

	const numGoroutines = 100
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			_, err := smsService.GetSMSMessage(context.Background(), uint(id))
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Unexpected error in race test: %v", err)
	}

	if atomic.LoadInt64(&accessCounter) != numGoroutines {
		t.Errorf("Expected %d SMS retrievals, got %d", numGoroutines, accessCounter)
	}
}

func TestSMSService_MixedOperations_RaceCondition(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var createCounter int64
	var retrieveCounter int64
	var balanceCounter int64

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		newID := atomic.AddInt64(&createCounter, 1)
		time.Sleep(1 * time.Millisecond)
		return domain.SMSID(newID), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		atomic.AddInt64(&balanceCounter, 1)
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		atomic.AddInt64(&retrieveCounter, 1)
		time.Sleep(1 * time.Millisecond)

		return &domain.SMS{
			ID:        filter.ID,
			Recipient: "+1234567890",
			Message:   "Mixed race test",
			Status:    string(domain.Delivered),
			CreatedAt: time.Now(),
		}, nil
	}

	const numOperations = 50
	var wg sync.WaitGroup
	errors := make(chan error, numOperations*2)

	for i := 0; i < numOperations; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()
			req := &presenter.SendSMSReq{
				UserID:    uint(id),
				Recipient: "+1234567890",
				Message:   "Mixed race test send",
			}
			_, err := smsService.SendSMS(context.Background(), req)
			if err != nil {
				errors <- err
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			_, err := smsService.GetSMSMessage(context.Background(), uint(id))
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Unexpected error in mixed race test: %v", err)
	}

	if atomic.LoadInt64(&createCounter) != numOperations {
		t.Errorf("Expected %d SMS creations, got %d", numOperations, createCounter)
	}

	if atomic.LoadInt64(&balanceCounter) != numOperations {
		t.Errorf("Expected %d balance updates, got %d", numOperations, balanceCounter)
	}

	if atomic.LoadInt64(&retrieveCounter) != numOperations {
		t.Errorf("Expected %d SMS retrievals, got %d", numOperations, retrieveCounter)
	}
}

func TestSMSService_SharedData_RaceCondition(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	sharedCounter := int64(0)
	var mutex sync.RWMutex

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		mutex.Lock()
		sharedCounter++
		currentValue := sharedCounter
		mutex.Unlock()

		time.Sleep(1 * time.Millisecond)
		return domain.SMSID(currentValue), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		mutex.RLock()
		_ = sharedCounter
		mutex.RUnlock()

		time.Sleep(1 * time.Millisecond)
		return nil
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		mutex.RLock()
		_ = sharedCounter
		mutex.RUnlock()

		time.Sleep(1 * time.Millisecond)

		return &domain.SMS{
			ID:        filter.ID,
			Recipient: "+1234567890",
			Message:   "Shared data test",
			Status:    string(domain.Delivered),
			CreatedAt: time.Now(),
		}, nil
	}

	const numGoroutines = 50
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*2)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()
			req := &presenter.SendSMSReq{
				UserID:    uint(id),
				Recipient: "+1234567890",
				Message:   "Shared data send test",
			}
			_, err := smsService.SendSMS(context.Background(), req)
			if err != nil {
				errors <- err
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			_, err := smsService.GetSMSMessage(context.Background(), uint(id+1))
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Unexpected error in shared data race test: %v", err)
	}

	mutex.RLock()
	finalCounter := sharedCounter
	mutex.RUnlock()

	if finalCounter != numGoroutines {
		t.Errorf("Expected final counter to be %d, got %d", numGoroutines, finalCounter)
	}
}

func TestSMSService_ContextCancellation_RaceCondition(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var completedOperations int64
	var cancelledOperations int64

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		select {
		case <-time.After(10 * time.Millisecond):
			atomic.AddInt64(&completedOperations, 1)
			return domain.SMSID(1), nil
		case <-ctx.Done():
			atomic.AddInt64(&cancelledOperations, 1)
			return 0, ctx.Err()
		}
	}

	const numGoroutines = 20
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			req := &presenter.SendSMSReq{
				UserID:    uint(id),
				Recipient: "+1234567890",
				Message:   "Context cancellation test",
			}

			_, _ = smsService.SendSMS(ctx, req)
		}(i)
	}

	time.Sleep(20 * time.Millisecond)
	cancel()
	wg.Wait()

	totalOperations := atomic.LoadInt64(&completedOperations) + atomic.LoadInt64(&cancelledOperations)
	if totalOperations == 0 {
		t.Error("Expected some operations to complete or be cancelled")
	}
}

func TestSMSService_ErrorHandling_RaceCondition(t *testing.T) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var successCount int64
	var errorCount int64

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		time.Sleep(1 * time.Millisecond)

		if recipient == "error" {
			atomic.AddInt64(&errorCount, 1)
			return 0, &MockError{message: "simulated error"}
		}

		atomic.AddInt64(&successCount, 1)
		return domain.SMSID(1), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	const numGoroutines = 100
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			recipient := "+1234567890"
			if id%10 == 0 {
				recipient = "error"
			}

			req := &presenter.SendSMSReq{
				UserID:    uint(id),
				Recipient: recipient,
				Message:   "Error handling race test",
			}

			_, _ = smsService.SendSMS(context.Background(), req)
		}(i)
	}

	wg.Wait()

	expectedErrors := int64(numGoroutines / 10)
	expectedSuccesses := int64(numGoroutines - expectedErrors)

	if atomic.LoadInt64(&errorCount) != expectedErrors {
		t.Errorf("Expected %d errors, got %d", expectedErrors, atomic.LoadInt64(&errorCount))
	}

	if atomic.LoadInt64(&successCount) != expectedSuccesses {
		t.Errorf("Expected %d successes, got %d", expectedSuccesses, atomic.LoadInt64(&successCount))
	}
}

func TestSMSService_HighLoad_RaceCondition(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high load race test in short mode")
	}

	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	var operationCount int64

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		newID := atomic.AddInt64(&operationCount, 1)
		return domain.SMSID(newID), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		atomic.AddInt64(&operationCount, 1)
		return &domain.SMS{
			ID:        filter.ID,
			Recipient: "+1234567890",
			Message:   "High load test",
			Status:    string(domain.Delivered),
			CreatedAt: time.Now(),
		}, nil
	}

	const numGoroutines = 1000
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				req := &presenter.SendSMSReq{
					UserID:    uint(id),
					Recipient: "+1234567890",
					Message:   "High load send test",
				}
				_, _ = smsService.SendSMS(context.Background(), req)
			} else {
				_, _ = smsService.GetSMSMessage(context.Background(), uint(id))
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("High load test completed in %v with %d operations", duration, atomic.LoadInt64(&operationCount))

	if duration > 10*time.Second {
		t.Errorf("High load test took too long: %v", duration)
	}
}
