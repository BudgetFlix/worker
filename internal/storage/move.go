package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"worker/internal/job"
)

func MoveJob(
	job *job.MediaJob,
	targetDir string,
) error {

	targetPath := filepath.Join(
		targetDir,
		filepath.Base(job.Path),
	)

	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf(
			"target job already exists: %s",
			targetPath,
		)
	}

	err := os.Rename(job.Path, targetPath)
	if err != nil {
		return err
	}

	job.Path = targetPath

	return nil
}