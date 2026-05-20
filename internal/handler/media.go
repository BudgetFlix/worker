package handler

import (
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/job"
	"worker/internal/logger"
	"worker/internal/pipeline"
	"worker/internal/storage"
)

var PipelineWG sync.WaitGroup

func Media(msg amqp.Delivery) error {

	logger.Loging("✅ Get message in handler")

	mediajob, err := job.FromJSON(msg.Body)
	if err != nil {
		return err
	}

	logger.Loging("✅ Success formating in handler")
	logger.Job(mediajob)

	err = job.Validate(mediajob)
	if err != nil {
		return err
	}

	logger.Loging("✅ Success validation in handler")

	for i := range mediajob.Items {
		err = storage.ChangeState(&mediajob.Items[i], storage.StateNew, storage.StateReady)
		if err != nil {
			return err
		}
	}

	logger.Loging("✅ Success add .ready")

	switch mediajob.Type {

	case job.MediaTypeMovie:
		pipe := pipeline.NewMovie(&PipelineWG)

		return pipe.Execute(
			context.Background(),
			mediajob,
		)

	default:
		return fmt.Errorf("unsuported media type %s", mediajob.Type)
	}

}
