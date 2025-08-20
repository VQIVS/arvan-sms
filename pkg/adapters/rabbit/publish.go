package rabbit

import (
	"fmt"
	"sms-dispatcher/pkg/constants"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func (r *Rabbit) Publish(Body []byte, Q string) error {
	pubErr := r.Ch.Publish(
		constants.Exchange,
		Q,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(Body),
		},
	)
	if pubErr != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v, message: %s", Q, pubErr, Body)
	}
	r.Logger.With("trace_id", uuid.NewString()).Info("published message to queue", "queue", Q, "message", string(Body))

	return nil

}
