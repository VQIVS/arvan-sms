package rabbit

import (
	"log/slog"
	"sms-dispatcher/pkg/constants"
	"sms-dispatcher/pkg/logger"
	"time"

	"github.com/streadway/amqp"
)

// TODO: do we need lock here?
func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		constants.Exchange,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			r.processMessage(d, queueName, handler)
		}
	}()

	return nil
}

func (r *Rabbit) processMessage(d amqp.Delivery, queueName string, handler func([]byte) error) {
	tracedLogger := logger.GetTracedLogger().With("queue", queueName)

	tracedLogger.Info("processing message from queue")

	if err := r.retryHandler(d.Body, handler, tracedLogger); err != nil {
		tracedLogger.Error("message processing failed after all retries", "error", err)
		r.nackMessage(d, tracedLogger)
		return
	}

	r.ackMessage(d, tracedLogger)
}

func (r *Rabbit) retryHandler(body []byte, handler func([]byte) error, logger *slog.Logger) error {
	const maxAttempts = 3
	const baseDelay = 100 * time.Millisecond

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := handler(body); err != nil {
			lastErr = err
			logger.Warn("message processing failed",
				"attempt", attempt,
				"max_attempts", maxAttempts,
				"error", err)

			if attempt < maxAttempts {
				delay := time.Duration(attempt) * baseDelay
				logger.Debug("retrying after delay", "delay", delay)
				time.Sleep(delay)
				continue
			}
		} else {
			// Success
			if attempt > 1 {
				logger.Info("message processed successfully after retries", "attempts", attempt)
			} else {
				logger.Info("message processed successfully")
			}
			return nil
		}
	}

	return lastErr
}

func (r *Rabbit) ackMessage(d amqp.Delivery, logger *slog.Logger) {
	if err := d.Ack(false); err != nil {
		logger.Error("failed to ack message", "error", err)
	}
}

func (r *Rabbit) nackMessage(d amqp.Delivery, logger *slog.Logger) {
	if err := d.Nack(false, false); err != nil {
		logger.Error("failed to nack message", "error", err)
	}
}
