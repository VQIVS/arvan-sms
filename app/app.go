package app

import (
	"context"
	"sms-dispatcher/config"
	"sms-dispatcher/internal/sms"
	"sms-dispatcher/internal/sms/port"
	"sms-dispatcher/pkg/adapters/rabbit"
	"sms-dispatcher/pkg/adapters/storage"
	appCtx "sms-dispatcher/pkg/context"
	"sms-dispatcher/pkg/postgres"

	"gorm.io/gorm"
)

type app struct {
	db         *gorm.DB
	cfg        config.Config
	smsService port.Service
	rabbit     *rabbit.Rabbit
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Rabbit() *rabbit.Rabbit {
	return a.rabbit
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) SMSService(ctx context.Context) port.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.smsService == nil {
			a.smsService = a.smsServiceWithDB(a.db)
		}
		return a.smsService
	}

	return a.smsServiceWithDB(db)
}
func (a *app) setDB() error {
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
	postgres.Migrate(db)

	a.db = db
	return nil
}

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}

	if err := a.setDB(); err != nil {
		return nil, err
	}
	// initialize rabbit connection if configured
	if cfg.Rabbit.URL != "" {
		r, err := rabbit.NewRabbit(cfg.Rabbit.URL)
		if err != nil {
			return nil, err
		}
		a.rabbit = r
		// initialize queues from config
		if err := a.rabbit.InitQueues(cfg.Rabbit.Queues); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func (a *app) smsServiceWithDB(db *gorm.DB) port.Service {
	return sms.NewService(storage.NewSMSRepo(db), a.rabbit)
}
