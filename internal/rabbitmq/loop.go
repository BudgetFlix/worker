package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


func Loop (
	msgs <-chan amqp.Delivery,
	handler func(amqp.Delivery) error,
){
	log.Println("consumer loop started")
	
	for msg := range msgs{

		err := handler(msg)

		if err != nil {
			log.Printf("message handle faild: %v", err)
			
			Nack(msg)

			continue
		}

		Ack(msg)

	}
	log.Println("consumer loop stopped")
}