# rubbit

Rubbit is a sample project i did to get more familiar with RabbitMQ and its plugins(delayed message).


Rubbit has the ability to publish messages in two different modes: normal and delayed.

It can consume messages as well and write them in redis.

In order to be more familiar with RabbitMq follow this link: `https://www.rabbitmq.com/tutorials/tutorial-one-go.html`


# installation
To get the Rubbit package use the command below:

go get github.com/bamdadam/rubbit

# Instructions

## Run

* run services and server: `make run-server`
* run services: `make up`
* run clients: `go run main.go client {topic names}`

# Structure

## Server
```go
func InitRabbitHandler() (*RabbitHandler, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("error while making connection to amqp broker: ", err)
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while making amqp channel: ", err)
		return nil, err
	}
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
```

## publish
```go
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
```

## publish delayed
```go
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
```

## client

```go
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
```

# Example

## Posting Message
* request
```
curl -X POST http://127.0.0.1:8080/publish -H 'Content-Type: application/json' -d '{"topic":"topic1", "message": "first-message", "publish_delay":"5000ms", "delayed":true}'
```

## Posting Delayed Message
* request
```
curl -X POST http://127.0.0.1:8080/publish -H 'Content-Type: application/json' -d '{"topic":"topic1", "message": "second-message", "delayed":false}'
```

## Getting Messages
* request
```
curl -X GET http://127.0.0.1:8080/subject -H 'Content-Type: application/json' -d '{"topic":"topic1"}'

```