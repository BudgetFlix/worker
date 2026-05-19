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
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

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

func MoveJobToAvailableDir(
	job *job.MediaJob,
	targetDir string,
) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	baseName := filepath.Base(job.Path)
	targetPath := filepath.Join(targetDir, baseName)

	for index := 1; ; index++ {
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			break
		} else if err != nil {
			return err
		}

		targetPath = filepath.Join(
			targetDir,
			fmt.Sprintf("%s_%d", baseName, index),
		)
	}

	err := os.Rename(job.Path, targetPath)
	if err != nil {
		return err
	}

	job.Path = targetPath

	return nil
}
