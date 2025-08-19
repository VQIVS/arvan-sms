package consumer

import (
	"context"
	"log/slog"
	"sms-dispatcher/app"
	"sms-dispatcher/pkg/constants"
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

func (h *Handler) Start(ctx context.Context) error {
	if h.app == nil || h.app.Rabbit() == nil {
		h.logger.Info("no rabbit configured, consumer won't start")
		return nil
	}

	svc := h.app.SMSService(context.Background())

	if err := h.app.Rabbit().Consume(constants.QueueSMSUpdate, func(body []byte) error {
		h.logger.Info("received message from queue", "queue", constants.QueueSMSUpdate, "message", string(body))
		return svc.UpdateSMSStatus(context.Background(), body)
	}); err != nil {
		h.logger.Error("failed to start consumer", "error", err)
		return err
	}
	<-ctx.Done()
	h.app.Rabbit().Close()
	h.logger.Info("consumer stopped")
	return nil
}
