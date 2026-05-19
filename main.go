package main

import (
	"log"

	"worker/internal/config"
	"worker/internal/rabbitmq"
	"worker/internal/handler"
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
		handler.Media,
	)
}