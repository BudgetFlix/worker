package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	channel, err := p.connection.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	if _, err := channel.QueueDeclarePassive(
		queue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("verify publish queue %s: %w", queue, err)
	}

	if err := channel.Confirm(false); err != nil {
		return fmt.Errorf("enable publisher confirms: %w", err)
	}

	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	returns := channel.NotifyReturn(make(chan amqp.Return, 1))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := channel.PublishWithContext(
		ctx,
		"",
		queue,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return err
	}

	return waitForPublishResult(queue, confirms, returns)
}

func waitForPublishResult(
	queue string,
	confirms <-chan amqp.Confirmation,
	returns <-chan amqp.Return,
) error {
	timeout := time.NewTimer(10 * time.Second)
	defer timeout.Stop()

	select {
	case returned := <-returns:
		return fmt.Errorf(
			"publish to %s was returned by broker: %s",
			queue,
			returned.ReplyText,
		)
	case confirmation, ok := <-confirms:
		if !ok {
			return fmt.Errorf("publisher confirmation channel closed")
		}
		if !confirmation.Ack {
			return fmt.Errorf("broker negatively acknowledged publish to %s", queue)
		}

		grace := time.NewTimer(100 * time.Millisecond)
		defer grace.Stop()
		select {
		case returned := <-returns:
			return fmt.Errorf(
				"publish to %s was returned by broker: %s",
				queue,
				returned.ReplyText,
			)
		case <-grace.C:
			return nil
		}
	case <-timeout.C:
		return fmt.Errorf("waiting for publish confirmation: timeout")
	}
}
