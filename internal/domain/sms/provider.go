package sms

import "context"

type SMSProvider interface {
	SendSMS(ctx context.Context, message *SMSMessage) (providerName string, err error)
}

type SMSProviderFunc func(ctx context.Context, message *SMSMessage) (string, error)

func (f SMSProviderFunc) SendSMS(ctx context.Context, message *SMSMessage) (string, error) {
	return f(ctx, message)
}
