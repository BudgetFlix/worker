package rabbitmq

import (
	"context"
	"log"
	"time"

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
				handleError(ctx, msg, err, producer, cfg)

				continue
			}
			handleDone(ctx, msg, producer, cfg)
		}
	}
}

func handleDone(
	ctx context.Context,
	msg amqp.Delivery,
	producer *Producer,
	cfg *config.Config,
) {

	mediaJob, err := job.FromJSON(msg.Body)
	if err != nil {
		log.Printf("failed to parse message: %v", err)
		if nackErr := Nack(msg); nackErr != nil {
			log.Printf("failed to nack malformed message: %v", nackErr)
		}
		return
	}

	retryMsg := dto.MediaRetryMsg{
		ID:     mediaJob.ID,
		Status: "DONE",
	}

	err = publishStatus(
		ctx,
		producer,
		cfg,
		cfg.UploadRetryQueue,
		retryMsg,
	)

	if err != nil {
		log.Printf("failed to publish retry message: %v", err)
		if nackErr := NackRequeue(msg); nackErr != nil {
			log.Printf("failed to nack message after publish failure: %v", nackErr)
		}
		return
	}

	if err := Ack(msg); err != nil {
		log.Printf("failed to ack completed message: %v", err)
		if nackErr := NackRequeue(msg); nackErr != nil {
			log.Printf("failed to requeue message after ack failure: %v", nackErr)
		}
	}
}

func handleError(
	ctx context.Context,
	msg amqp.Delivery,
	processErr error,
	producer *Producer,
	cfg *config.Config,
) {

	mediaJob, err := job.FromJSON(msg.Body)
	if err != nil {
		log.Printf("failed to parse message: %v", err)
		if nackErr := Nack(msg); nackErr != nil {
			log.Printf("failed to nack malformed message: %v", nackErr)
		}
		return
	}

	retryMsg := dto.MediaRetryMsg{
		ID:       mediaJob.ID,
		Status:   "ERROR",
		ErrorMsg: processErr.Error(),
	}

	err = publishStatus(
		ctx,
		producer,
		cfg,
		cfg.UploadRetryQueue,
		retryMsg,
	)
	if err != nil {
		log.Printf("failed to publish retry message: %v", err)
		if nackErr := NackRequeue(msg); nackErr != nil {
			log.Printf("failed to nack message after publish failure: %v", nackErr)
		}
		return
	}

	if err := Nack(msg); err != nil {
		log.Printf("failed to nack failed message: %v", err)
		if requeueErr := NackRequeue(msg); requeueErr != nil {
			log.Printf("failed to requeue message after nack failure: %v", requeueErr)
		}
	}
}

func publishStatus(
	ctx context.Context,
	producer *Producer,
	cfg *config.Config,
	queue string,
	message dto.MediaRetryMsg,
) error {
	persisted := false

	for {
		if err := producer.Publish(queue, message); err == nil {
			if persisted {
				if removeErr := RemoveStatus(cfg, message); removeErr != nil {
					log.Printf("failed to remove delivered status outbox entry: %v", removeErr)
				}
			}
			return nil
		} else {
			log.Printf("failed to publish status, retrying: %v", err)
		}

		if !persisted {
			if err := PersistStatus(cfg, queue, message); err != nil {
				log.Printf("failed to persist status outbox entry: %v", err)
			} else {
				persisted = true
			}
		}

		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
