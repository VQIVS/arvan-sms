package config

type Config struct {
	Server   Server   `yaml:"server"`
	Database DBConfig `yaml:"database"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RabbitMQ struct {
	Host     string  `yaml:"host"`
	Port     int     `yaml:"port"`
	Username string  `yaml:"username"`
	Password string  `yaml:"password"`
	Queues   []Queue `yaml:"queues"`
}

type Queue struct {
	Name     string `yaml:"name"`
	Exchange string `yaml:"exchange"`
	Routing  string `yaml:"routing"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
