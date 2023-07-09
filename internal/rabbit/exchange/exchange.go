package exchange

import "github.com/rabbitmq/amqp091-go"

func NewTopicExchange(ch *amqp091.Channel, name string) error {
	err := ch.ExchangeDeclare(name, "x-delayed-message", false, false, false, false, amqp091.Table{
		"x-delayed-type": "topic",
	})
	return err
}
