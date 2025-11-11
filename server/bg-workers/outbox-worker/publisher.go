package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(OutboxEvent) error
	Close()
}

type RabbitPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewRabbitPublisher() (*RabbitPublisher, error) {
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

	return &RabbitPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *RabbitPublisher) Publish(evt OutboxEvent) error {
	// Convert OutboxEvent to JSON
	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	// Publish to the queue
	return p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent, // durable messages
			Timestamp:    time.Now(),
			Type:         evt.EventType,
		},
	)
}

func (p *RabbitPublisher) Close() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
