package storage

import (
	"os"
	"path/filepath"
	"testing"

	"worker/internal/config"
)

func TestRecoverProcessingJobsMovesJobsToErrorDir(t *testing.T) {
	root := t.TempDir()
	processDir := filepath.Join(root, "process")
	errorDir := filepath.Join(root, "error")
	jobDir := filepath.Join(processDir, "job_123")

	if err := os.MkdirAll(jobDir, 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{
		"movie.mp4.processing",
		"trailer.mp4.ready",
		"poster.jpg",
	}

	for _, file := range files {
		if err := os.WriteFile(filepath.Join(jobDir, file), []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	cfg := config.Config{
		ProcessDir: processDir,
		ErrorDir:   errorDir,
	}

	if err := RecoverProcessingJobs(cfg); err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"movie.mp4.error",
		"trailer.mp4.error",
		"poster.jpg.error",
	} {
		if _, err := os.Stat(filepath.Join(errorDir, "job_123", file)); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := os.Stat(jobDir); !os.IsNotExist(err) {
		t.Fatalf("expected processing job dir to be moved, got error: %v", err)
	}
}

func TestRecoverProcessingJobsIgnoresMissingProcessDir(t *testing.T) {
	cfg := config.Config{
		ProcessDir: filepath.Join(t.TempDir(), "missing"),
		ErrorDir:   filepath.Join(t.TempDir(), "error"),
	}

	if err := RecoverProcessingJobs(cfg); err != nil {
		t.Fatal(err)
	}
}
