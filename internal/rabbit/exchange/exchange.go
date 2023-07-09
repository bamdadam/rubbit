package exchange

import "github.com/rabbitmq/amqp091-go"

func NewDelayedTopicExchange(ch *amqp091.Channel, name string) error {
	err := ch.ExchangeDeclare(name, "x-delayed-message", false, true, false, false, amqp091.Table{
		"x-delayed-type": "topic",
	})
	return err
}

func NewTopicExchange(ch *amqp091.Channel, name string) error {
	err := ch.ExchangeDeclare(name, "topic", false, true, false, false, nil)
	return err
}
