package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *Connection
}

func NewConsumer (connection *Connection) *Consumer{

	return &Consumer{
		connection: connection,
	}
}

func (c *Consumer) Consume( queue string) (<- chan amqp.Delivery,error){
	return c.connection.channel.Consume(queue,"",false,false,false,false,nil)
}

