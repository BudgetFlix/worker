package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/config"
	"worker/internal/rabbitmq"
)

func main() {

	cfg := config.Load()

	connection, err := rabbitmq.NewConnection(
		cfg.RabbitMQURL(),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	consumer := rabbitmq.NewConsumer(
		connection,
	)

	msgs, err := consumer.Consume(
		"video.upload.queue",
	)

	if err != nil {
		log.Fatal(err)
	}

	rabbitmq.Loop(
		msgs,
		handleMessage,
	)
}

func handleMessage(
	msg amqp.Delivery,
) error {

	log.Printf(
		"processing message: %s",
		string(msg.Body),
	)

	return nil
}