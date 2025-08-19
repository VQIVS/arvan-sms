package http

import (
	"context"
	"sms-dispatcher/api/service"
)

type SMSServiceGetter interface {
	GetSMSService(ctx context.Context) *service.SMSService
}
