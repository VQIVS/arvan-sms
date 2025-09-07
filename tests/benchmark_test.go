package tests

import (
	"context"
	"sms/internal/domain/sms"
	"testing"
	"time"
)


func BenchmarkSMSMessage_MarkAsSent(b *testing.B) {
	message := &sms.SMSMessage{
		ID:       "test-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}
	provider := "test-provider"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message.Status = sms.SMSStatusPending
		message.MarkAsSent(provider)
	}
}

func BenchmarkSMSMessage_MarkAsFailed(b *testing.B) {
	message := &sms.SMSMessage{
		ID:       "test-id",
		UserID:   "user-123",
		Content:  "Test message",
		Receiver: "+1234567890",
		Status:   sms.SMSStatusPending,
	}
	provider := "test-provider"
	failureCode := "NETWORK_ERROR"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message.Status = sms.SMSStatusPending
		message.MarkAsFailed(provider, failureCode)
	}
}

func BenchmarkEventCreation(b *testing.B) {
	userID := "user-123"
	smsID := "sms-456"
	amount := int64(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := sms.RequestSMSBilling{
			UserID:    userID,
			SMSID:     smsID,
			Amount:    amount,
			TimeStamp: time.Now(),
		}
		_ = event
	}
}

func BenchmarkMockRepo_Create(b *testing.B) {
	repo := newMockSMSRepo()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message := &sms.SMSMessage{
			ID:       generateBenchmarkID(i),
			UserID:   "user-123",
			Content:  "Test message",
			Receiver: "+1234567890",
			Status:   sms.SMSStatusPending,
		}
		repo.Create(ctx, message)
	}
}

func BenchmarkMockRepo_GetByFilter(b *testing.B) {
	repo := newMockSMSRepo()
	ctx := context.Background()

	for i := 0; i < 1000; i++ {
		message := &sms.SMSMessage{
			ID:       generateBenchmarkID(i),
			UserID:   "user-123",
			Content:  "Test message",
			Receiver: "+1234567890",
			Status:   sms.SMSStatusPending,
		}
		repo.Create(ctx, message)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := generateBenchmarkID(i % 1000)
		filter := sms.Filter{ID: &id}
		_, _ = repo.GetByFilter(ctx, filter)
	}
}

func BenchmarkPhoneNumberValidation(b *testing.B) {
	phoneNumbers := []string{
		"+1234567890",
		"+989123456789",
		"+441234567890",
		"+861234567890",
		"+33123456789",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		phone := phoneNumbers[i%len(phoneNumbers)]
		isValidE164PhoneNumber(phone)
	}
}

func BenchmarkSMSContentValidation(b *testing.B) {
	contents := []string{
		"Short message",
		"This is a medium length message for testing purposes",
		generateString(160),
		"Hello 😊 with emojis",
		"Numbers: 123456789",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content := contents[i%len(contents)]
		isValidSMSContent(content)
	}
}

func generateBenchmarkID(i int) string {
	return "benchmark-id-" + string(rune('0'+i%10)) + string(rune('a'+i%26))
}

func BenchmarkSMSMessageAllocation(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		message := &sms.SMSMessage{
			ID:       "test-id",
			UserID:   "user-123",
			Content:  "Test message",
			Receiver: "+1234567890",
			Status:   sms.SMSStatusPending,
		}
		_ = message
	}
}

func BenchmarkSMSMessage_MarkAsSent_Parallel(b *testing.B) {
	provider := "test-provider"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			message := &sms.SMSMessage{
				ID:       "test-id",
				UserID:   "user-123",
				Content:  "Test message",
				Receiver: "+1234567890",
				Status:   sms.SMSStatusPending,
			}
			message.MarkAsSent(provider)
		}
	})
}
