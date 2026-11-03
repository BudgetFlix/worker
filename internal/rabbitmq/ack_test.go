package rabbitmq

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestAckReturnsBrokerError(t *testing.T) {
	if err := Ack(amqp.Delivery{}); err == nil {
		t.Fatal("expected ack error from delivery without a channel")
	}
}

func TestNackReturnsBrokerError(t *testing.T) {
	if err := Nack(amqp.Delivery{}); err == nil {
		t.Fatal("expected nack error from delivery without a channel")
	}
}

func TestNackRequeueReturnsBrokerError(t *testing.T) {
	if err := NackRequeue(amqp.Delivery{}); err == nil {
		t.Fatal("expected requeue error from delivery without a channel")
	}
}
