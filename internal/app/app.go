package app

import (
	"context"
	"sms/config"
	"sms/internal/usecase/sms"
	"sms/pkg/rabbit"

	"gorm.io/gorm"
)

type app struct {
	db         *gorm.DB
	cfg        config.Config
	rabbitConn *rabbit.RabbitConn
	smsService *sms.UseCase
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) RabbitConn() *rabbit.RabbitConn {
	return a.rabbitConn
}

func (a *app) SMSService(ctx context.Context) *sms.UseCase {
	return a.smsService
}
