package rabbit

import (
	"github.com/streadway/amqp"
)

func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		"amq.topic",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	const maxAttempts = 3

	go func() {
		for d := range msgs {
			r.Logger.Info("processing message from queue", "queue", queueName)

			if err := handler(d.Body); err != nil {
				attempts := 0
				if d.Headers != nil {
					if v, ok := d.Headers["attempts"]; ok {
						switch t := v.(type) {
						case int:
							attempts = t
						case int32:
							attempts = int(t)
						case int64:
							attempts = int(t)
						case float32:
							attempts = int(t)
						case float64:
							attempts = int(t)
						}
					}
				}

				attempts++

				if attempts <= maxAttempts {
					retryQ := queueName + ".retry"
					headers := amqp.Table{}
					if d.Headers != nil {
						for k, v := range d.Headers {
							headers[k] = v
						}
					}
					headers["attempts"] = attempts

					pubErr := r.Ch.Publish(
						"",
						retryQ,
						false,
						false,
						amqp.Publishing{
							Headers:      headers,
							ContentType:  d.ContentType,
							Body:         d.Body,
							DeliveryMode: amqp.Persistent,
						},
					)
					if pubErr != nil {
						r.Logger.Error("failed to publish to retry queue", "queue", retryQ, "error", pubErr)
						_ = d.Nack(false, true)
						continue
					}

					if err := d.Ack(false); err != nil {
						r.Logger.Error("failed to ack original message after republish", "error", err)
					} else {
						r.Logger.Info("republished message to retry queue", "queue", retryQ, "attempt", attempts)
					}
					continue
				}

				dlq := queueName + ".dlq"
				headers := amqp.Table{}
				if d.Headers != nil {
					for k, v := range d.Headers {
						headers[k] = v
					}
				}
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
					continue
				}

				if err := d.Ack(false); err != nil {
					r.Logger.Error("failed to ack original message after dlq publish", "error", err)
				} else {
					r.Logger.Info("moved message to DLQ", "queue", dlq, "attempt", attempts)
				}

			} else {
				if err := d.Ack(false); err != nil {
					r.Logger.Error("failed to ack message", "queue", queueName, "error", err)
				} else {
					r.Logger.Info("successfully processed message from queue", "queue", queueName)
				}
			}
		}
	}()

	return nil
}
