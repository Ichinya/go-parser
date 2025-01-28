package services

import (
	"log"

	"github.com/rabbitmq/amqp091-go" // Используем amqp091-go
)

func ConnectToRabbitMQ(uri string) (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, err
	}

	log.Println("Успешно подключено к RabbitMQ!")
	return conn, nil
}
