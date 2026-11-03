package pipeline

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"worker/internal/config"
	"worker/internal/job"
)

func TestMoveFailedJobMovesProcessingJobToErrorDir(t *testing.T) {
	root := t.TempDir()
	processDir := filepath.Join(root, "process")
	errorDir := filepath.Join(root, "error")
	jobPath := filepath.Join(processDir, "job_123")
	processingPath := filepath.Join(jobPath, "movie.mp4.processing")

	if err := os.MkdirAll(jobPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(processingPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	mediaJob := &job.MediaJob{
		ID:   "123",
		Path: jobPath,
		Items: []job.VideoItem{{
			FileName: "movie.mp4.processing",
			Path:     processingPath,
		}},
	}
	cause := errors.New("artwork move failed")
	cfg := config.Config{ErrorDir: errorDir}

	if err := moveFailedJob(cfg, mediaJob, &mediaJob.Items[0], cause); !errors.Is(err, cause) {
		t.Fatalf("expected original error, got %v", err)
	}

	errorJobPath := filepath.Join(errorDir, "job_123")
	errorFilePath := filepath.Join(errorJobPath, "movie.mp4.error")
	if _, err := os.Stat(errorFilePath); err != nil {
		t.Fatalf("expected failed job file at %s: %v", errorFilePath, err)
	}
	if _, err := os.Stat(jobPath); !os.IsNotExist(err) {
		t.Fatalf("expected processing job to be moved, got %v", err)
	}
}
