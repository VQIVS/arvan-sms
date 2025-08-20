package rabbit

import (
	"github.com/streadway/amqp"
)

func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		"",
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
			r.Logger.Info("processing message from queue", "queue", queueName)

			if err := r.processMessageWithRetry(queueName, d, handler); err != nil {
				r.Logger.Error("message processing failed after all retries", "queue", queueName, "error", err)
			}
		}
	}()

	return nil
}

func (r *Rabbit) processMessageWithRetry(queueName string, d amqp.Delivery, handler func([]byte) error) error {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := handler(d.Body); err != nil {
			r.Logger.Error("message processing failed", "queue", queueName, "attempt", attempt, "error", err)

			if attempt == maxAttempts {
				// Send to DLQ after all retries failed
				r.sendToDLQ(queueName, d, attempt)
				return err
			}

			r.Logger.Info("retrying message processing", "queue", queueName, "attempt", attempt+1)
			continue
		}

		// Success
		if err := d.Ack(false); err != nil {
			r.Logger.Error("failed to ack message", "queue", queueName, "error", err)
			return err
		}

		r.Logger.Info("successfully processed message from queue", "queue", queueName, "attempt", attempt)
		return nil
	}

	return nil
}

func (r *Rabbit) sendToDLQ(queueName string, d amqp.Delivery, attempts int) {
	dlq := queueName + ".dlq"
	headers := r.copyHeaders(d.Headers)
	headers["attempts"] = attempts

	if err := r.Ch.Publish(
		"",
		dlq,
		false,
		false,
		amqp.Publishing{
			Headers:      headers,
			ContentType:  d.ContentType,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		r.Logger.Error("failed to publish to dlq", "queue", dlq, "error", err)
		_ = d.Nack(false, true)
		return
	}

	if err := d.Ack(false); err != nil {
		r.Logger.Error("failed to ack original message after dlq publish", "error", err)
	} else {
		r.Logger.Info("moved message to DLQ", "queue", dlq, "attempt", attempts)
	}
}

func (r *Rabbit) copyHeaders(original amqp.Table) amqp.Table {
	headers := amqp.Table{}
	for k, v := range original {
		headers[k] = v
	}
	return headers
}
