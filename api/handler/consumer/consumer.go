package consumer

import (
	"log/slog"
	"sms-dispatcher/app"
	"sms-dispatcher/pkg/logger"
)

type Handler struct {
	app    app.App
	logger *slog.Logger
}

func New(a app.App) *Handler {
	return &Handler{
		app:    a,
		logger: logger.GetLogger(),
	}
}
