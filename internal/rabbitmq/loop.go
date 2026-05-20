package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Loop(
	ctx context.Context,
	msgs <-chan amqp.Delivery,
	handler func(amqp.Delivery) error,
) {

	log.Println("consumer loop started")

	for {
		select {

		case <-ctx.Done():
			log.Println("consumer loop stopped")
			return

		case msg, ok := <-msgs:

			if !ok {
				log.Println("consumer channel closed")
				return
			}

			err := handler(msg)

			if err != nil {

				log.Printf(
					"message handle failed: %v",
					err,
				)

				Nack(msg)

				continue
			}

			Ack(msg)
		}
	}
}