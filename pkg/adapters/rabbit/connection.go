package rabbit

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitProvider struct {
	uri      string
	exchange string
	conn     *amqp.Connection
	ch       *amqp.Channel
	confirms chan amqp.Confirmation
}

func NewRabbitWithConn(uri, exchange string) (*RabbitProvider, error) {
	r := &RabbitProvider{
		uri:      uri,
		exchange: exchange,
	}
	if err := r.connect(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RabbitProvider) connect() error {
	c, err := amqp.Dial(r.uri)
	if err != nil {
		return err
	}
	ch, err := c.Channel()
	if err != nil {
		return err
	}

	if err = ch.ExchangeDeclarePassive(r.exchange, "topic", true, false, false, false, nil); err != nil {
		r.Close()
		return err
	}
	err = ch.Confirm(false)
	if err != nil {
		r.Close()
		return err
	}

	r.confirms = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	r.conn = c
	r.ch = ch
	return nil

}

func (r *RabbitProvider) DeclareAndBind(q QueueSpec) error {
	if err := r.validateQueueSpec(q); err != nil {
		return err
	}
	if err := r.declareQueue(q.Name, q.Args); err != nil {
		return err
	}
	if err := r.bindQueueToExchange(q.Name, q.Bindings); err != nil {
		return err
	}
	return nil
}

func (r *RabbitProvider) Close() error {
	err := r.ch.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitProvider) validateQueueSpec(q QueueSpec) error {
	if q.Name == "" {
		return errors.New("queue name required")
	}
	if q.Args == nil {
		q.Args = amqp.Table{}
	}
	return nil
}

func (r *RabbitProvider) declareQueue(name string, args amqp.Table) error {
	_, err := r.ch.QueueDeclare(name, true, false, false, false, args)
	if err != nil {
		return fmt.Errorf("queue declare %s: %w", name, err)
	}
	return nil
}

func (r *RabbitProvider) bindQueueToExchange(queue string, bindings []string) error {
	for _, key := range bindings {
		if err := r.ch.QueueBind(queue, key, r.exchange, false, nil); err != nil {
			return fmt.Errorf("bind %s -> %s: %w", key, queue, err)
		}
	}
	return nil
}
