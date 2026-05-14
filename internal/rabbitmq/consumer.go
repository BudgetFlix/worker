package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *Consumer) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}