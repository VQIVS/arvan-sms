package rabbit

import (
	"log/slog"
	"sms-dispatcher/pkg/logger"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	Conn   *amqp.Connection
	Ch     *amqp.Channel
	Logger *slog.Logger
}

func NewRabbit(url string, customLogger *slog.Logger) (*Rabbit, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	log := customLogger
	if log == nil {
		log = logger.GetLogger()
	}

	return &Rabbit{Conn: conn, Ch: ch, Logger: log}, nil
}

func (r *Rabbit) Close() {
	if r.Ch != nil {
		_ = r.Ch.Close()
		_ = r.Conn.Close()
	}
}

func (r *Rabbit) InitQueues(keys []string, exchange string) error {
	for _, key := range keys {
		if err := r.declareBind(exchange, key, true); err != nil {
			r.Logger.Error("failed to declare and bind queue", "key", key, "error", err)
			return err
		}
		r.Logger.Info("queue declared and bound", "key", key)
	}
	return nil
}

func (r *Rabbit) declareBind(exchange string, routingKey string, durable bool) error {
	if r == nil || r.Ch == nil {
		return nil
	}
	q, err := r.Ch.QueueDeclare(
		GetQueueName(routingKey),
		durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.Ch.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
