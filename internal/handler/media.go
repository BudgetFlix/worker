package handler

import(
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/job"
	"worker/internal/pipeline"
)

func Media (msg amqp.Delivery) error {

	mediajob, err := job.FromJSON(msg.Body,) 

	if err != nil {return err}

	err = job.Validate(mediajob)
	if err != nil {
		return err
	} 

	switch mediajob.Type {

		case job.MediaTypeMovie:
			pipe := pipeline.NewMovie()

			return pipe.Execute(
				context.Background(),
				mediajob,
			)
		
		default:
			return fmt.Errorf("unsuported media type %s",mediajob.Type)
	}

}

