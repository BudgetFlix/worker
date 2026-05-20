package handler

import(
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/job"
	"worker/internal/pipeline"
	"worker/internal/logger"
)

func Media (msg amqp.Delivery) error {

	logger.Loging("✅ Get message in handler")
	
	mediajob, err := job.FromJSON(msg.Body,) 
	if err != nil {return err}
	
	logger.Loging("✅ Sucsess formating in handler")
	logger.Job(mediajob)
	
	err = job.Validate(mediajob)
	if err != nil {
		return err
	} 
	
	logger.Loging("✅ Sucsess validation in handler")
		
	//todo: set extension to .ready

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

