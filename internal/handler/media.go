package handler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/config"
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

	cfg := config.Load()
	if filepath.Clean(filepath.Dir(mediajob.Path)) == filepath.Clean(cfg.DoneDir) {
		return nil
	}

	logger.Loging("✅ Success formating in handler")
	logger.Job(mediajob)

	err = job.Validate(mediajob)
	if err != nil {
		return failMediaJob(mediajob, err)
	}

	logger.Loging("✅ Success validation in handler")

	switch mediajob.Type {
	case job.MediaTypeMovie:
		if len(mediajob.Items) != 1 {
			return failMediaJob(
				mediajob,
				errors.New("movie pipeline requires exactly one video item"),
			)
		}
	default:
		return failMediaJob(
			mediajob,
			fmt.Errorf("unsuported media type %s", mediajob.Type),
		)
	}

	for i := range mediajob.Items {
		item := &mediajob.Items[i]
		fromState := storage.StateNew
		if strings.HasSuffix(item.Path, string(storage.StateError)) {
			fromState = storage.StateError
		}

		err = storage.ChangeState(item, fromState, storage.StateReady)
		if err != nil {
			return failMediaJob(mediajob, err)
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

func failMediaJob(mediajob *job.MediaJob, cause error) error {
	cfg := config.Load()
	errs := []error{cause}

	for i := range mediajob.Items {
		if err := storage.ChangeToError(&mediajob.Items[i]); err != nil {
			errs = append(errs, fmt.Errorf("mark failed item: %w", err))
		}
	}

	if filepath.Clean(filepath.Dir(mediajob.Path)) != filepath.Clean(cfg.ErrorDir) {
		if err := storage.MoveJob(mediajob, cfg.ErrorDir); err != nil {
			errs = append(errs, fmt.Errorf("move failed job to error directory: %w", err))
		}
	}

	return errors.Join(errs...)
}
