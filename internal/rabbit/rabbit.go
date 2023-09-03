package rabbit

import (
	"context"
	"log"
	"time"

	"github.com/bamdadam/rubbit/internal/rabbit/exchange"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitHandler struct {
	ch     *amqp091.Channel
	con    *amqp091.Connection
	ctx    context.Context
	cancel context.CancelFunc
}

func InitRabbitHandler() (*RabbitHandler, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("error while making connection to amqp broker: ", err)
		return nil, err
	}
	// defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while making amqp channel: ", err)
		return nil, err
	}
	// defer ch.Close()
	err = exchange.NewDelayedTopicExchange(ch, "Announcer-Delayed")
	if err != nil {
		log.Println("error while making amqp exchange: ", err)
		return nil, err
	}

	err = exchange.NewTopicExchange(ch, "Announcer")
	if err != nil {
		log.Println("error while making amqp exchange: ", err)
		return nil, err
	}
	rb := new(RabbitHandler)
	rb.ch = ch
	rb.con = conn
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rb.ctx = ctx
	rb.cancel = cancel
	return rb, nil
}
