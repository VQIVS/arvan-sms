package rabbit

import (
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
)

func (r *RabbitProvider) Publish(routingKey string, body []byte) error {
	if r.ch == nil {
		return errors.New("no channel")
	}

	err := r.ch.Publish(
		r.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish %s: %w", routingKey, err)
	}
	select {
	case c := <-r.confirms:
		if !c.Ack {
			return errors.New("publish not acknowledged")
		}
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("publish confirm timeout within 5 seconds")
	}
}

func (r *RabbitProvider) StartSubscribe(q QueueSpec, handler func(amqp.Delivery)) (<-chan amqp.Delivery, error) {
	prefetch := q.Prefetch
	if prefetch <= 0 {
		prefetch = 32
	}

	if err := r.ch.Qos(prefetch, 0, false); err != nil {
		return nil, err
	}

	tag := q.ConsumerTag
	if tag == "" {
		tag = q.Name + "_c"
	}

	deliveries, err := r.ch.Consume(
		q.Name,
		tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	go func() {
		for d := range deliveries {
			retryOperation := func() error {
				handler(d)
				return nil
			}

			backOff := backoff.NewExponentialBackOff()
			backOff.InitialInterval = 1 * time.Second
			backOff.MaxInterval = 30 * time.Second
			backOff.MaxElapsedTime = 5 * time.Minute

			err := backoff.Retry(retryOperation, backOff)
			if err != nil {

				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return deliveries, nil
}
