package app

import (
	"context"
	"sms/config"
	"sms/internal/infra/external"
	"sms/internal/infra/messaging"
	"sms/internal/infra/storage"
	"sms/internal/infra/storage/types"
	"sms/internal/usecase/sms"
	"sms/pkg/postgres"
	"sms/pkg/rabbit"

	"gorm.io/gorm"
)

type app struct {
	db         *gorm.DB
	cfg        config.Config
	rabbitConn *rabbit.RabbitConn
	smsService *sms.Service
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

func (a *app) SMSService(ctx context.Context) *sms.Service {
	return a.smsService
}

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}
	if err := a.setDB(); err != nil {
		return nil, err
	}

	if err := a.setRabbitConn(); err != nil {
		return nil, err
	}

	if err := a.initQueues(); err != nil {
		return nil, err
	}

	a.smsService = setService(a.db, a.rabbitConn)
	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func setService(db *gorm.DB, rabbitConn *rabbit.RabbitConn) *sms.Service {
	smsRepo := storage.NewSMSRepository(db)
	smsPublisher := messaging.NewSMSPublisher(rabbitConn)
	smsProvider := external.DefaultSMSProvider()
	return sms.NewSMSService(smsRepo, smsPublisher, smsProvider, db)
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
	// Auto migrate
	err = postgres.Migrate(db, &types.SMS{})
	if err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *app) setRabbitConn() error {
	rabbitConn := rabbit.NewRabbitConn(a.cfg.RabbitMQ.URI)
	a.rabbitConn = rabbitConn
	return nil
}

func (a *app) initQueues() error {
	for _, q := range a.cfg.RabbitMQ.Queues {
		err := a.rabbitConn.DeclareBindQueue(q.Name, q.Exchange, q.Routing)
		if err != nil {
			return err
		}
	}
	return nil
}
