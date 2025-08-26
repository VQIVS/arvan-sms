package rabbit

import (
	"github.com/streadway/amqp"
)

type QueueSpec struct {
	Name        string
	Bindings    []string
	Prefetch    int
	Args        amqp.Table
	ConsumerTag string
}
