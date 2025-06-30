package config

import (
	"log"
	"github.com/streadway/amqp"
)

func InitRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("falied to connect to RabbitMQ: %s", err)
	}
	return conn
}