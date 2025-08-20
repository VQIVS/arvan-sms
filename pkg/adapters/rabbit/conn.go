package rabbit

import (
	"log/slog"
	"sms-dispatcher/pkg/constants"
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
		conn.Close()
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

func (r *Rabbit) InitQueues(queue string) error {
	if r == nil || r.Ch == nil {
		return nil
	}
	queueName := GetQueueName(queue)
	_, err := r.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	r.Ch.QueueBind(queueName, constants.KeySMSUpdate, "amq.topic", false, nil)
	if err != nil {
		return err
	}
	return err
}

func GetQueueName(key string) string {
	return "sms_" + key
}
