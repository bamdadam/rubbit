package rabbit

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bamdadam/rubbit/internal/rabbit/queue"
	"github.com/bamdadam/rubbit/internal/store/rdb"
	"github.com/rabbitmq/amqp091-go"
)

func InitRabbitClient(rdb *rdb.RedisStore) error {
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
	q, err := queue.NewQueue(ch)
	if err != nil {
		log.Println("error while making amqp queue: ", err)
		return err
	}
	args := getArgs()
	for _, arg := range args {
		err = ch.QueueBind(q.Name, arg, "Announcer", false, nil)
		if err != nil {
			log.Println("error while binding queue to amqp exchange: ", err)
			return err
		}
		err = ch.QueueBind(q.Name, arg, "Announcer-Delayed", false, nil)
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
			err := handleMessage(msg, rdb)
			if err != nil {
				log.Print("can't handle message: ", err)
			}
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
	return nil
}

func handleMessage(message amqp091.Delivery, rdb *rdb.RedisStore) error {
	log.Println("message is: ", string(message.Body))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := rdb.SaveMessage(ctx, string(message.Body), message.Timestamp.Format("2006-01-02 15:04:05"), message.RoutingKey)
	return err
}

func getArgs() []string {
	args := os.Args
	if len(args) <= 2 {
		args = append(args, "#")
	}
	return args[2:]
}
