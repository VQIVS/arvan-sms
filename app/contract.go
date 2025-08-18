package app

import (
	"context"
	"sms-dispatcher/config"
	"sms-dispatcher/internal/sms/port"
	"sms-dispatcher/pkg/adapters/rabbit"

	"gorm.io/gorm"
)

type App interface {
	SMSService(ctx context.Context) port.Service
	DB() *gorm.DB
	Config() config.Config
	Rabbit() *rabbit.Rabbit
}
