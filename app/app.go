package app

import (
	"context"
	"log/slog"
	"sms-dispatcher/config"
	"sms-dispatcher/internal/sms"
	"sms-dispatcher/internal/sms/port"
	"sms-dispatcher/pkg/adapters/rabbit"
	"sms-dispatcher/pkg/adapters/storage"
	"sms-dispatcher/pkg/constants"
	"sms-dispatcher/pkg/logger"
	"sms-dispatcher/pkg/postgres"

	"gorm.io/gorm"
)

type App struct {
	db         *gorm.DB
	cfg        config.Config
	smsService port.Service
	rabbit     *rabbit.Rabbit
	logger     *slog.Logger
}

func (a *App) SMSService(ctx context.Context) port.Service {
	if a.smsService == nil {
		a.smsService = sms.NewService(storage.NewSMSRepo(a.db), a.rabbit)
	}
	return a.smsService
}

func (a *App) setDB() error {
	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   a.cfg.DB.User,
		Pass:   a.cfg.DB.Password,
		Host:   a.cfg.DB.Host,
		Port:   a.cfg.DB.Port,
		DBName: a.cfg.DB.Database,
		Schema: a.cfg.DB.Schema,
	})

	if err != nil {
		return err
	}
	// auto migrate gorm
	if err := postgres.Migrate(db); err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *App) setRabbit() error {
	rabbit, err := rabbit.NewRabbitWithConn(
		a.cfg.Rabbit.URL,
		constants.TopicExchange,
	)
	if err != nil {
		return err
	}
	a.rabbit = rabbit
	return nil
}

func NewApp(cfg config.Config) (*App, error) {
	l := logger.GetLogger()
	a := &App{
		cfg:    cfg,
		logger: l,
	}
	if err := a.setDB(); err != nil {
		return nil, err
	}
	if err := a.setRabbit(); err != nil {
		return nil, err
	}
	return a, nil
}

func NewMustApp(cfg config.Config) *App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}
