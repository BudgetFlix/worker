package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/rabbitmq"

	"worker/internal/config"

)

func main() {
	cfg := config.Load()


	consumer, err := rabbitmq.NewConsumer(
		cfg.RabbitMQURL(),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	msgs, err := consumer.Consume("video.upload.queue")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("worker started")

	for msg := range msgs {
		err := handleMessage(msg)

		if err != nil {
			log.Println(err)

			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}
}

func handleMessage(msg amqp.Delivery) error {
	log.Printf("processing: %s", string(msg.Body))

	return nil
}