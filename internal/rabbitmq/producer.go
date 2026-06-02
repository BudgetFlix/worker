package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connection *Connection
}

func NewProducer(
	connection *Connection,
) *Producer {

	return &Producer{
		connection: connection,
	}
}

func (p *Producer) Publish(
	queue string,
	message any,
) error {

	body, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return p.connection.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
