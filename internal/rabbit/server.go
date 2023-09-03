package rabbit

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitHandler) PublishDelayedMessage(topic string, message string, pubDelay int64) error {
	err := r.ch.PublishWithContext(
		r.ctx,
		"Announcer-Delayed",
		topic,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Headers: map[string]interface{}{
				"x-delay": pubDelay,
			},
			Timestamp: time.Now().Add(time.Duration(pubDelay) * time.Millisecond),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitHandler) PublishMessage(topic string, message string) error {
	err := r.ch.PublishWithContext(
		r.ctx,
		"Announcer",
		topic,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
