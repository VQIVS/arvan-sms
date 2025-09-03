package app

import (
	"context"
	"sms/config"
	"sms/internal/usecase/sms"
	"sms/pkg/rabbit"

	"gorm.io/gorm"
)

type App interface {
	Config() config.Config
	DB() *gorm.DB
	RabbitConn() *rabbit.RabbitConn
	SMSService(ctx context.Context) *sms.UseCase
}
