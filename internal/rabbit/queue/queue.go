package queue

import "github.com/rabbitmq/amqp091-go"

func NewQueue(ch *amqp091.Channel) (amqp091.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
