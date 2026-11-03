package rabbitmq

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestWaitForPublishResultRejectsNegativeConfirmation(t *testing.T) {
	confirms := make(chan amqp.Confirmation, 1)
	returns := make(chan amqp.Return, 1)
	confirms <- amqp.Confirmation{Ack: false}

	if err := waitForPublishResult("status", confirms, returns); err == nil {
		t.Fatal("expected negative confirmation error")
	}
}

func TestWaitForPublishResultRejectsReturnedMessage(t *testing.T) {
	confirms := make(chan amqp.Confirmation, 1)
	returns := make(chan amqp.Return, 1)
	returns <- amqp.Return{ReplyText: "unroutable"}

	if err := waitForPublishResult("status", confirms, returns); err == nil {
		t.Fatal("expected returned message error")
	}
}

func TestWaitForPublishResultDetectsLateReturn(t *testing.T) {
	confirms := make(chan amqp.Confirmation, 1)
	returns := make(chan amqp.Return, 1)
	confirms <- amqp.Confirmation{Ack: true}

	go func() {
		time.Sleep(10 * time.Millisecond)
		returns <- amqp.Return{ReplyText: "unroutable"}
	}()

	if err := waitForPublishResult("status", confirms, returns); err == nil {
		t.Fatal("expected late returned message error")
	}
}
