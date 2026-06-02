package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/config"
	"worker/internal/dto"
	"worker/internal/job"
)

func Loop(
	ctx context.Context,
	msgs <-chan amqp.Delivery,
	handler func(amqp.Delivery) error,
	producer *Producer,
	cfg *config.Config,
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
				handleError(msg, err, producer, cfg)

				continue
			}
			handleDone(msg, producer, cfg)
		}
	}
}

func handleDone(
	msg amqp.Delivery,
	producer *Producer,
	cfg *config.Config,
) {

	mediaJob, err := job.FromJSON(msg.Body)
	if err != nil {
		log.Printf("failed to parse message: %v", err)
		Nack(msg)
		return
	}

	retryMsg := dto.MediaRetryMsg{
		ID:     mediaJob.ID,
		Status: "DONE",
	}

	err = producer.Publish(
		cfg.UploadRetryQueue,
		retryMsg,
	)

	if err != nil {
		log.Printf("failed to publish retry message: %v", err)
		Nack(msg)
		return
	}

	Ack(msg)
}

func handleError(
	msg amqp.Delivery,
	processErr error,
	producer *Producer,
	cfg *config.Config,
) {

	mediaJob, err := job.FromJSON(msg.Body)
	if err != nil {
		log.Printf("failed to parse message: %v", err)
		Nack(msg)
		return
	}

	retryMsg := dto.MediaRetryMsg{
		ID:       mediaJob.ID,
		Status:   "ERROR",
		ErrorMsg: processErr.Error(),
	}

	err = producer.Publish(
		cfg.UploadRetryQueue,
		retryMsg,
	)
	if err != nil {
		log.Printf("failed to publish retry message: %v", err)
		Nack(msg)
		return
	}

	Nack(msg)
}
