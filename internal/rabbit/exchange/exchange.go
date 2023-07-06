package exchange

import "github.com/rabbitmq/amqp091-go"

func NewTopicExchange(ch *amqp091.Channel, name string) error {
	err := ch.ExchangeDeclare(name, "topic", false, false, false, false, nil)
	return err
}
