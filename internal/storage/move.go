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
    return moveArtwork(
        mediaJob,
        movieSource,
        "poster.jpg",
    )
}

func MoveBackground(
    mediaJob *job.MediaJob,
    movieSource string,
) error {
    return moveArtwork(
        mediaJob,
        movieSource,
        "background.jpg",
    )
}

func moveArtwork(
    mediaJob *job.MediaJob,
    movieSource string,
    filename string,
) error {
    sourcePath := filepath.Join(
        mediaJob.Path,
        filename,
    )

    if _, err := os.Stat(sourcePath); err != nil {
        if os.IsNotExist(err) {
            return nil
        }

        return fmt.Errorf(
            "stat %s: %w",
            filename,
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
        filename,
    )

    if _, err := os.Stat(targetPath); err == nil {
        return fmt.Errorf(
            "target %s already exists: %s",
            filename,
            targetPath,
        )
    } else if !os.IsNotExist(err) {
        return fmt.Errorf(
            "stat target %s: %w",
            filename,
            err,
        )
    }

    if err := os.Rename(sourcePath, targetPath); err != nil {
        return fmt.Errorf(
            "move %s: %w",
            filename,
            err,
        )
    }

    return nil
}