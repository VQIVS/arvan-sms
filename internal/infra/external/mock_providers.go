package external

import (
	"context"
	"errors"
	"math/rand"
	"sms/internal/domain/sms"
	"time"
)

const (
	MockProviderName   = "MockProvider"
	RandomFailProvider = "RandomFailProvider"
	AlwaysFailProvider = "AlwaysFailProvider"
)

func MockSMSProvider() sms.SMSProviderFunc {
	return func(ctx context.Context, message *sms.SMSMessage) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return MockProviderName, nil
	}
}

func RandomFailSMSProvider(failProbability float64) sms.SMSProviderFunc {
	return func(ctx context.Context, message *sms.SMSMessage) (string, error) {
		time.Sleep(100 * time.Millisecond)

		if rand.Float64() < failProbability {
			return RandomFailProvider, errors.New("random delivery failure")
		}

		return RandomFailProvider, nil
	}
}

func AlwaysFailSMSProvider() sms.SMSProviderFunc {
	return func(ctx context.Context, message *sms.SMSMessage) (string, error) {
		return AlwaysFailProvider, errors.New("delivery failure")
	}
}
