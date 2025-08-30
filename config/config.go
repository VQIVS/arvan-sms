package config

type Config struct {
	Server   Server
	RabbitMQ RabbitMQ
}

type Server struct {
	Host string
	Port int
}

type RabbitMQ struct {
	Host     string
	Port     int
	Username string
	Password string
	Queues   []Queue
}

type Queue struct {
	Name     string
	Exchange string
	Routing  string
}
