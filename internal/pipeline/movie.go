package pipeline

import (
	"context"
	"errors"
	"fmt"

	"worker/internal/config"
	"worker/internal/ffmpeg"
	"worker/internal/job"
	"worker/internal/storage"
	"worker/internal/logger"
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

	logger.Loging("The basic: ")
	logger.Job(job)
	
	err := storage.MoveJob(
		job,
		cfg.ProcessDir,
	)
	
	if err != nil {
		return err
	}
	logger.Loging("✅ Succses moving to process dir ")
	logger.Job(job)
	
	item := &job.Items[0]
	input := job.ItemPath(item)

	
	processingPath, err := storage.ChangeState(
		input,
		storage.StateNew,
		storage.StateProcessing,
	)
	
	if err != nil {
		return moveFailedJob(cfg, job, "", err)
	}
	logger.Loging("✅ succses set state to .process")
	logger.Job(job)
	
	outputDir, err := storage.CreateMovieLibraryDir(
		cfg,
		job,
	)
	
	if err != nil {
		return moveFailedJob(cfg, job, processingPath, err)
	}
	logger.Loging("✅ create the outdir")
	logger.Job(job)
	
	
	logger.Loging("Start the encode")
	err = ffmpeg.HLS(
		ctx,
		processingPath,
		outputDir,
	)
	
	
	if err != nil {
		return moveFailedJob(cfg, job, processingPath, err)
	}
	logger.Loging("✅ succses encode")
	logger.Job(job)
	
	_, err = storage.ChangeState(
		processingPath,
		storage.StateProcessing,
		storage.StateDone,
	)
	
	
	if err != nil {
		return moveFailedJob(cfg, job, processingPath, err)
	}
	logger.Loging("✅ succses .done extension")
	logger.Job(job)
	
	err = storage.MoveJobToAvailableDir(
		job,
		cfg.DoneDir,
	)
	
	if err != nil {
		return err
	}
	logger.Loging("✅ succses moveing to done dir")
	logger.Job(job)

	return nil
}

func moveFailedJob(
	cfg config.Config,
	job *job.MediaJob,
	processingPath string,
	cause error,
) error {
	var stateErr error
	if processingPath != "" {
		_, stateErr = storage.ChangeState(
			processingPath,
			storage.StateProcessing,
			storage.StateError,
		)
	}

	moveErr := storage.MoveJobToAvailableDir(
		job,
		cfg.ErrorDir,
	)

	errs := []error{cause}
	if stateErr != nil {
		errs = append(errs, fmt.Errorf("mark failed job: %w", stateErr))
	}

	if moveErr != nil {
		errs = append(errs, fmt.Errorf("move failed job to error directory: %w", moveErr))
	}

	if len(errs) > 1 {
		return errors.Join(errs...)
	}

	return cause
}
