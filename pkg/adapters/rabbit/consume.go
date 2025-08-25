package rabbit

import (
	"log/slog"
	"sms-dispatcher/pkg/constants"
	"time"

	backOff "github.com/cenkalti/backoff/v4"

	"github.com/streadway/amqp"
)

func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		constants.TopicExchange,
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
			if err := handler(d.Body); err != nil {
				r.nackMessage(d, r.Logger)
				r.Logger.Error("failed to handle message", "error", err)
				continue
			}
			r.ackMessage(d, r.Logger)
		}
	}()

	return nil
}

func (r *Rabbit) ConsumeWithRetry(queueName string, handler func([]byte) error, maxElapsedTime time.Duration) error {
	operation := func() error {
		return r.Consume(queueName, handler)
	}
	b := backOff.NewExponentialBackOff()
	b.MaxElapsedTime = maxElapsedTime
	return backOff.Retry(operation, b)
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
