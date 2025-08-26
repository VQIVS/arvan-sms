package rabbit

import (
	"github.com/streadway/amqp"
)

// Rabbit represents a connection to a RabbitMQ server, including channel and exchange information.
type Rabbit struct {
	uri      string                 // Connection URI for RabbitMQ
	exchange string                 // Exchange name (e.g., amq.topic)
	conn     *amqp.Connection       // AMQP connection
	ch       *amqp.Channel          // AMQP channel
	confirms chan amqp.Confirmation // Channel for publisher confirms
}

// QueueSpec defines the specification for a RabbitMQ queue and its consumer settings.
type QueueSpec struct {
	Name        string     // Name of the queue
	Bindings    []string   // List of routing keys to bind to the queue
	Prefetch    int        // Prefetch count for consumer (messages to fetch before ack)
	Args        amqp.Table // Additional arguments for queue declaration
	ConsumerTag string     // Consumer tag for identifying the consumer
}
