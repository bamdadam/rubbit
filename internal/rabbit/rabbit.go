package rabbit

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bamdadam/rubbit/internal/rabbit/exchange"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitServer struct {
	ch     *amqp091.Channel
	con    *amqp091.Connection
	ctx    context.Context
	cancel context.CancelFunc
}

func InitRabbitServer() (*RabbitServer, error) {
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
	err = exchange.NewTopicExchange(ch, "Announcer")
	if err != nil {
		log.Println("error while making amqp exchange: ", err)
		return nil, err
	}
	rb := new(RabbitServer)
	rb.ch = ch
	rb.con = conn
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rb.ctx = ctx
	rb.cancel = cancel
	return rb, nil
}

func (r *RabbitServer) PublishMessage(topic string, message string, pubDelay int64) error {
	err := r.ch.PublishWithContext(
		r.ctx,
		"Announcer",
		topic,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Headers: map[string]interface{}{
				"x-delay": pubDelay,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func InitRabbitClient() error {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("error while making connection to amqp broker")
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while making amqp channel: ", err)
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Println("error while making amqp queue: ", err)
		return err
	}
	args := os.Args
	if len(args) <= 2 {
		args = append(args, "#")
	}
	for _, arg := range args[2:] {
		err = ch.QueueBind(q.Name, arg, "Announcer", false, nil)
		if err != nil {
			log.Println("error while binding queue to amqp exchange: ", err)
			return err
		}
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("error while registering a consumer: ", err)
		return err
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			log.Println("message is: ", string(msg.Body))
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
	return nil
}
