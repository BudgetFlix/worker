package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"worker/internal/job"
)

func MoveJob(
	job *job.MediaJob,
	targetDir string,
) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	oldPath := job.Path

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

	if err := os.Rename(job.Path, targetPath); err != nil {
		return err
	}

	job.Path = targetPath

	for i := range job.Items {
		item := &job.Items[i]

		item.Path = strings.Replace(
			item.Path,
			oldPath,
			job.Path,
			1,
		)
	}

	return nil
}