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

func BenchmarkSMSService_SendSMS(b *testing.B) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Benchmark test message",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := smsService.SendSMS(context.Background(), req)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkSMSService_GetSMSMessage(b *testing.B) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSMS := &domain.SMS{
		ID:        domain.SMSID(123),
		Recipient: "+1234567890",
		Message:   "Benchmark test message",
		Status:    string(domain.Delivered),
		CreatedAt: time.Now(),
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		return mockSMS, nil
	}

	smsID := uint(123)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := smsService.GetSMSMessage(context.Background(), smsID)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkSMSService_SendSMS_Parallel(b *testing.B) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   "Parallel benchmark test message",
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := smsService.SendSMS(context.Background(), req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})
}

func BenchmarkSMSService_GetSMSMessage_Parallel(b *testing.B) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSMS := &domain.SMS{
		ID:        domain.SMSID(123),
		Recipient: "+1234567890",
		Message:   "Parallel benchmark test message",
		Status:    string(domain.Delivered),
		CreatedAt: time.Now(),
	}

	mockSvc.GetSMSByFilterFunc = func(ctx context.Context, filter *domain.SMSFilter) (*domain.SMS, error) {
		return mockSMS, nil
	}

	smsID := uint(123)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := smsService.GetSMSMessage(context.Background(), smsID)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})
}

func BenchmarkSMSService_SendSMS_LargeMessage(b *testing.B) {
	mockSvc := &MockService{}
	smsService := service.NewSMSService(mockSvc)

	mockSvc.CreateSMSFunc = func(ctx context.Context, recipient string, message string) (domain.SMSID, error) {
		return domain.SMSID(123), nil
	}

	mockSvc.UserBalanceUpdateFunc = func(ctx context.Context, user event.UserBalanceEvent) error {
		return nil
	}

	largeMessage := ""
	for i := 0; i < 1000; i++ {
		largeMessage += "A"
	}

	req := &presenter.SendSMSReq{
		UserID:    1,
		Recipient: "+1234567890",
		Message:   largeMessage,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := smsService.SendSMS(context.Background(), req)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
