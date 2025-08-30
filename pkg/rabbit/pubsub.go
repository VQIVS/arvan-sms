package rabbit

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Publisher struct {
	rabbitConn *RabbitConn
}

type Consumer struct {
	rabbitConn *RabbitConn
	handlers   map[string]func([]byte) error
}

func NewPublisher(conn *RabbitConn) *Publisher {
	return &Publisher{
		rabbitConn: conn,
	}
}

func (p *Publisher) Publish(queueName, exchange string, body interface{}) error {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return p.rabbitConn.Ch.Publish(
		exchange,
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJson,
		},
	)
}

func NewConsumer(conn *RabbitConn) *Consumer {
	return &Consumer{
		rabbitConn: conn,
		handlers:   make(map[string]func([]byte) error),
	}
}

func (c *Consumer) Subscribe(queueName string, handler func([]byte) error) {
	c.handlers[queueName] = handler
}

func (c *Consumer) StartConsume() error {
	for queueName, handler := range c.handlers {
		go c.consumeFromQueue(queueName, handler)
	}
	return nil
}

func (c *Consumer) consumeFromQueue(queueName string, handler func([]byte) error) {
	msgs, err := c.rabbitConn.Ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consuming from %s: %v", queueName, err)
		return
	}
	log.Printf("Consumer started for queue: %s", queueName)

	for msg := range msgs {
		if err := handler(msg.Body); err != nil {
			log.Printf("Error handling message in %s: %v", queueName, err)
			msg.Nack(false, true)
		} else {
			msg.Ack(false)
		}
	}
}

func (c *Consumer) SetQos(prefetchCount int) error {
	return c.rabbitConn.Ch.Qos(
		prefetchCount,
		0,
		false,
	)
}
