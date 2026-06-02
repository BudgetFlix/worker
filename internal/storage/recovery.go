package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"worker/internal/config"
)

func RecoverProcessingJobs(
	cfg config.Config,
) error {
	entries, err := os.ReadDir(cfg.ProcessDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err := os.MkdirAll(cfg.ErrorDir, 0755); err != nil {
		return err
	}

	var errs []error

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		jobPath := filepath.Join(
			cfg.ProcessDir,
			entry.Name(),
		)

		if err := recoverProcessingJob(jobPath, cfg.ErrorDir); err != nil {
			errs = append(
				errs,
				err,
			)
		}
	}

	return errors.Join(errs...)
}

func recoverProcessingJob(
	jobPath string,
	errorDir string,
) error {
	if err := markJobFilesError(jobPath); err != nil {
		return err
	}

	targetPath := filepath.Join(
		errorDir,
		filepath.Base(jobPath),
	)

	if err := os.Rename(jobPath, targetPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf(
			"move recovered job to error directory: %w",
			err,
		)
	}

	return nil
}

func markJobFilesError(
	jobPath string,
) error {
	return filepath.WalkDir(
		jobPath,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}

				return err
			}

			if entry.IsDir() {
				return nil
			}

			return markFileError(path)
		},
	)
}

func markFileError(
	path string,
) error {
	if strings.HasSuffix(path, string(StateError)) {
		return nil
	}

	targetPath := errorFilePath(path)

	if err := os.Rename(path, targetPath); err != nil {
		return fmt.Errorf(
			"mark recovered file as error: %w",
			err,
		)
	}

	return nil
}

func errorFilePath(
	path string,
) string {
	for _, state := range []FileState{
		StateProcessing,
		StateReady,
		StateDone,
	} {
		if strings.HasSuffix(path, string(state)) {
			return strings.TrimSuffix(
				path,
				string(state),
			) + string(StateError)
		}
	}

	return path + string(StateError)
}
