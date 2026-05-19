package storage

import (
	"os"
	"path/filepath"
	"strconv"

	"worker/internal/config"
	"worker/internal/job"
)

func CreateMovieLibraryDir(
	cfg config.Config,
	job *job.MediaJob,
) (string, error) {

	path := movieLibraryPath(
		cfg,
		job.MediaID,
	)

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

func movieLibraryPath(
	cfg config.Config,
	mediaID int64,
) string {
	return filepath.Join(
		cfg.MovieSource,
		strconv.FormatInt(mediaID, 10),
		"hls",
	)
}
