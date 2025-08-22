package rabbit

import (
	"fmt"
	"sms-dispatcher/pkg/logger"

	"github.com/streadway/amqp"
)

func (r *Rabbit) Publish(Body []byte, exchange string, routingKey string) error {
	err := r.Ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(Body),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v, message: %s", routingKey, err, Body)
	}
	logger.GetTracedLogger().Info("published message to queue", "queue", routingKey, "message", string(Body))

	return nil

}
