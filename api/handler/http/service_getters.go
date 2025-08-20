package http

import (
	"context"
	"sms-dispatcher/api/service"
	"sms-dispatcher/app"
	"sms-dispatcher/config"
)

type smsServiceProvider struct {
	appContainer app.App
	cfg          config.ServerConfig
}

func (p *smsServiceProvider) GetSMSService(ctx context.Context) *service.SMSService {
	return service.NewSMSService(p.appContainer.SMSService(ctx))
}

func newSMSServiceGetter(appContainer app.App, cfg config.ServerConfig) SMSServiceGetter {
	return &smsServiceProvider{
		appContainer: appContainer,
		cfg:          cfg,
	}
}
