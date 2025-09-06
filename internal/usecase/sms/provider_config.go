package sms

import (
	"context"
	"sms/internal/domain/sms"
	"sms/internal/infra/external"
)

func (u *Service) WithMockProvider() *Service {
	u.provider = external.MockSMSProvider()
	return u
}

func (u *Service) WithRandomFailProvider(failProbability float64) *Service {
	u.provider = external.RandomFailSMSProvider(failProbability)
	return u
}

func (u *Service) WithAlwaysFailProvider() *Service {
	u.provider = external.AlwaysFailSMSProvider()
	return u
}

func (u *Service) WithCustomProvider(provider sms.SMSProvider) *Service {
	u.provider = provider
	return u
}

func (u *Service) WithCustomProviderFunc(fn func(ctx context.Context, message *sms.SMSMessage) (string, error)) *Service {
	u.provider = sms.SMSProviderFunc(fn)
	return u
}
