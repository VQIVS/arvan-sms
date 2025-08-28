package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"sms-dispatcher/app"
	"sms-dispatcher/internal/sms/port"
	"sms-dispatcher/pkg/adapters/rabbit"
	"sms-dispatcher/pkg/constants"
	"sms-dispatcher/pkg/logger"

	"github.com/streadway/amqp"
)

type Handler struct {
	app        app.App
	logger     *slog.Logger
	smsService port.Service

	mu       sync.RWMutex
	stopChan chan struct{}
	stopped  bool
}

func NewHandler(a app.App) *Handler {
	return &Handler{
		app:        a,
		logger:     logger.GetLogger(),
		smsService: a.SMSService(context.Background()),
		stopChan:   make(chan struct{}),
	}
}

func (h *Handler) Start(ctx context.Context) error {
	if h.app == nil {
		return fmt.Errorf("app is nil")
	}

	if h.app.Rabbit() == nil {
		h.logger.Info("no rabbit configured, consumer won't start")
		return nil
	}

	if h.smsService == nil {
		return fmt.Errorf("SMS service is nil")
	}

	h.logger.Info("starting SMS consumer")

	queueName := rabbit.GetQueueName(constants.KeySMSUpdate)
	// TODO: queue names in config or somewhere else
	if err := h.app.Rabbit().InitQueues([]string{constants.KeySMSUpdate}); err != nil {
		h.logger.Error("failed to initialize queues", "error", err)
		return fmt.Errorf("failed to initialize queues: %w", err)
	}

	_, err := h.app.Rabbit().Subscribe(queueName, h.createMessageHandler())
	if err != nil {
		h.logger.Error("failed to start consumer", "error", err, "queue", queueName)
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	h.logger.Info("consumer started successfully", "queue", queueName)

	return h.handleGracefulShutdown(ctx)
}

func (h *Handler) createMessageHandler() func(amqp.Delivery) error {
	return func(delivery amqp.Delivery) error {
		startTime := time.Now()
		h.logger.Info("processing message",
			"queue", rabbit.GetQueueName(constants.KeySMSUpdate),
			"body_size", len(delivery.Body),
		)

		err := h.smsService.UpdateSMSStatus(context.Background(), delivery.Body)
		if err != nil {
			// TODO: Publish to DLQ then consume and refund the user OR publish to refund queue
			h.logger.Error("failed to process message",
				"error", err,
				"processing_time", time.Since(startTime),
			)
			return err
		}
		h.logger.Debug("message processed successfully",
			"processing_time", time.Since(startTime),
		)
		return nil
	}
}

// TODO: handle drain messages in the queue before shutdown
func (h *Handler) handleGracefulShutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		h.logger.Info("context cancelled, shutting down consumer")
	case <-h.stopChan:
		h.logger.Info("stop signal received, shutting down consumer")
	}

	return h.shutdown()
}

func (h *Handler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.stopped {
		return nil
	}

	h.logger.Info("stopping consumer")
	close(h.stopChan)
	h.stopped = true

	return nil
}

// cleanup
func (h *Handler) shutdown() error {
	h.logger.Info("shutting down consumer")

	if h.app.Rabbit() != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			defer close(done)
			h.app.Rabbit().Close()
		}()

		select {
		case <-done:
			h.logger.Info("rabbit connection closed successfully")
		case <-shutdownCtx.Done():
			h.logger.Warn("rabbit connection close timed out")
		}
	}

	h.logger.Info("consumer shutdown complete")
	return nil
}
