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

func MovePoster(
    mediaJob *job.MediaJob,
    movieSource string,
) error {
    posterPath := filepath.Join(
        mediaJob.Path,
        "poster.jpg",
    )

    if _, err := os.Stat(posterPath); err != nil {
        if os.IsNotExist(err) {
            return nil
        }

        return fmt.Errorf(
            "stat poster: %w",
            err,
        )
    }

    movieDir := filepath.Join(
        movieSource,
        fmt.Sprintf("%d", mediaJob.MediaID),
    )

    if err := os.MkdirAll(movieDir, 0755); err != nil {
        return fmt.Errorf(
            "create movie directory: %w",
            err,
        )
    }

    targetPath := filepath.Join(
        movieDir,
        "poster.jpg",
    )

    if _, err := os.Stat(targetPath); err == nil {
        return fmt.Errorf(
            "target poster already exists: %s",
            targetPath,
        )
    } else if !os.IsNotExist(err) {
        return fmt.Errorf(
            "stat target poster: %w",
            err,
        )
    }

    if err := os.Rename(posterPath, targetPath); err != nil {
        return fmt.Errorf(
            "move poster: %w",
            err,
        )
    }

    return nil
}