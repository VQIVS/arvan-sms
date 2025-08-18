package http

import (
	"context"
	"sms-dispatcher/api/service"
	"sms-dispatcher/app"
	"sms-dispatcher/config"
)

func smsServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.SMSService] {
	return func(ctx context.Context) *service.SMSService {
		return service.NewSMSService(appContainer.SMSService(ctx))
	}
}
