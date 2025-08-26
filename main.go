package main

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
)

func main() {
	q := amqp.QueueConfig{
		GenerateName: func(topic string) string {
			return "sms" + "key"
		},
		Durable: true,
	}
	conn, err := amqp.NewConnection(amqp.ConnectionConfig{
		AmqpURI: "amqp://guest:guest@localhost:5672/",
	}, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println(conn, q)
	defer conn.Close()
}
