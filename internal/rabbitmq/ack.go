package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Ack(msg amqp.Delivery) error {
	return msg.Ack(false)
}

func Nack(msg amqp.Delivery) error {
	return msg.Nack(false, false)
}

func NackRequeue(msg amqp.Delivery) error {
	return msg.Nack(false, true)
}
