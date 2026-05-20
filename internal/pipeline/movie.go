package pipeline

import (
	"context"
	"errors"
	"fmt"

	"worker/internal/config"
	"worker/internal/ffmpeg"
	"worker/internal/job"
	"worker/internal/logger"
	"worker/internal/storage"
)

type Movie struct{}

func NewMovie() *Movie {
	return &Movie{}
}

func (p *Movie) Execute(
	ctx context.Context,
	job *job.MediaJob,
) error {

	if len(job.Items) != 1 {
		return errors.New(
			"movie pipeline requires exactly one video item",
		)
	}

	cfg := config.Load()

	logger.Loging("The basic:")
	logger.Job(job)

	err := storage.MoveJob(
		job,
		cfg.ProcessDir,
	)

	if err != nil {
		return err
	}

	logger.Loging("✅ Success moving to process dir")
	logger.Job(job)

	item := &job.Items[0]

	err = storage.ChangeState(
		item,
		storage.StateNew,
		storage.StateProcessing,
	)

	if err != nil {
		return moveFailedJob(cfg, job, item, err)
	}

	logger.Loging("✅ Success set state to .processing")
	logger.Job(job)

	job.Outdir , err = storage.CreateMovieLibraryDir(
		cfg,
		job,
	)

	if err != nil {
		return moveFailedJob(cfg, job, item, err)
	}

	logger.Loging("✅ Create the outdir")
	logger.Job(job)

	logger.Loging("Start the encode")

	err = ffmpeg.HLS(
		ctx,
		item.Path,
		job.Outdir,
	)

	if err != nil {
		return moveFailedJob(cfg, job, item, err)
	}

	logger.Loging("✅ Success encode")
	logger.Job(job)

	err = storage.ChangeState(
		item,
		storage.StateProcessing,
		storage.StateDone,
	)

	if err != nil {
		return moveFailedJob(cfg, job, item, err)
	}

	logger.Loging("✅ Success .done extension")
	logger.Job(job)

	err = storage.MoveJob(
		job,
		cfg.DoneDir,
	)

	if err != nil {
		return err
	}

	logger.Loging("✅ Success moving to done dir")
	logger.Job(job)

	return nil
}

func moveFailedJob(
	cfg config.Config,
	job *job.MediaJob,
	item *job.VideoItem,
	cause error,
) error {

	var stateErr error

	if item != nil {
		stateErr = storage.ChangeState(
			item,
			storage.StateProcessing,
			storage.StateError,
		)
	}

	moveErr := storage.MoveJob(
		job,
		cfg.ErrorDir,
	)

	errs := []error{cause}

	if stateErr != nil {
		errs = append(
			errs,
			fmt.Errorf(
				"mark failed job: %w",
				stateErr,
			),
		)
	}

	if moveErr != nil {
		errs = append(
			errs,
			fmt.Errorf(
				"move failed job to error directory: %w",
				moveErr,
			),
		)
	}

	if len(errs) > 1 {
		return errors.Join(errs...)
	}

	return cause
}