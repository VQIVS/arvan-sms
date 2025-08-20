package consumer

import (
	"context"
	"log/slog"
	"sms-dispatcher/app"
	"sms-dispatcher/pkg/adapters/rabbit"
	"sms-dispatcher/pkg/constants"
	"sms-dispatcher/pkg/logger"

	"github.com/google/uuid"
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
	queue := rabbit.GetQueueName(constants.KeySMSUpdate)

	if err := h.app.Rabbit().Consume(queue, func(body []byte) error {
		h.logger.With("trace_id", uuid.NewString()).Info("received message from queue", "queue", queue, "message", string(body))
		return svc.UpdateSMSStatus(context.Background(), body)
	}); err != nil {
		h.logger.With("trace_id", uuid.NewString()).Error("failed to start consumer", "error", err)
		return err
	}
	<-ctx.Done()
	// h.app.Rabbit().Close()
	h.logger.With("trace_id", uuid.NewString()).Info("consumer stopped")
	return nil
}
