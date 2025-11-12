package initializers

import (
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	Queue   amqp091.Queue
}

func NewRabbitConsumer() (*RabbitConsumer, error) {
	url := os.Getenv("MQ_URL")
	qname := os.Getenv("MQ_QUEUE")
	if qname == "" {
		qname = "cliprelay.events"
	}

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Durable queue for event replay guarantee
	q, err := ch.QueueDeclare(
		qname,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	return &RabbitConsumer{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil
}
