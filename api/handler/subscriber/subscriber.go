package subscriber

import (
	"log/slog"
	"sms-dispatcher/app"
	"sms-dispatcher/pkg/logger"
)

type Subscriber struct {
	app    app.App
	logger *slog.Logger
}

func NewSubscriber(a app.App) *Subscriber {
	return &Subscriber{
		app:    a,
		logger: logger.GetLogger(),
	}
}

func (s *Subscriber) StartHandler() {
}
