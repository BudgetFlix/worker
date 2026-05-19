package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"worker/internal/job"
)

func TestMovieExecuteMovesJobToErrorWhenReadyStateIsMissing(t *testing.T) {
	dirs := setPipelineEnv(t)
	jobDir := filepath.Join(dirs.newDir, "job_1234")
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		t.Fatalf("create job dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(jobDir, "movie.mp4"), []byte("video"), 0644); err != nil {
		t.Fatalf("create video file: %v", err)
	}

	mediaJob := &job.MediaJob{
		ID:      "1234",
		MediaID: 42,
		Type:    job.MediaTypeMovie,
		Path:    jobDir,
		Items: []job.VideoItem{
			{Index: 0, FileName: "movie.mp4", State: job.StateReceived},
		},
		State: job.StateReceived,
	}

	err := NewMovie().Execute(context.Background(), mediaJob)
	if err == nil {
		t.Fatal("expected movie pipeline to fail")
	}

	expectedPath := filepath.Join(dirs.errorDir, "job_1234")
	if mediaJob.Path != expectedPath {
		t.Fatalf("expected job path %q, got %q", expectedPath, mediaJob.Path)
	}

	if _, err := os.Stat(filepath.Join(expectedPath, "movie.mp4")); err != nil {
		t.Fatalf("expected job file in error dir: %v", err)
	}
}

func TestMovieExecuteMovesProcessingJobToError(t *testing.T) {
	dirs := setPipelineEnv(t)
	jobDir := filepath.Join(dirs.newDir, "job_1234")
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		t.Fatalf("create job dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(jobDir, "movie.mp4.ready"), []byte("video"), 0644); err != nil {
		t.Fatalf("create video file: %v", err)
	}

	blockingFile := filepath.Join(t.TempDir(), "movies-file")
	if err := os.WriteFile(blockingFile, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("create blocking movie source file: %v", err)
	}
	t.Setenv("MOVIE_SOURCE", blockingFile)

	mediaJob := &job.MediaJob{
		ID:      "1234",
		MediaID: 42,
		Type:    job.MediaTypeMovie,
		Path:    jobDir,
		Items: []job.VideoItem{
			{Index: 0, FileName: "movie.mp4.ready", State: job.StateReceived},
		},
		State: job.StateReceived,
	}

	err := NewMovie().Execute(context.Background(), mediaJob)
	if err == nil {
		t.Fatal("expected movie pipeline to fail")
	}

	expectedPath := filepath.Join(dirs.errorDir, "job_1234")
	if mediaJob.Path != expectedPath {
		t.Fatalf("expected job path %q, got %q", expectedPath, mediaJob.Path)
	}

	if _, err := os.Stat(filepath.Join(expectedPath, "movie.mp4.error")); err != nil {
		t.Fatalf("expected error-state file in error dir: %v", err)
	}
}

type pipelineDirs struct {
	newDir     string
	processDir string
	doneDir    string
	errorDir   string
}

func setPipelineEnv(t *testing.T) pipelineDirs {
	t.Helper()

	root := t.TempDir()
	dirs := pipelineDirs{
		newDir:     filepath.Join(root, "new"),
		processDir: filepath.Join(root, "process"),
		doneDir:    filepath.Join(root, "done"),
		errorDir:   filepath.Join(root, "error"),
	}

	t.Setenv("RABBITMQ_HOST", "localhost")
	t.Setenv("RABBITMQ_PORT", "5672")
	t.Setenv("RABBITMQ_USERNAME", "guest")
	t.Setenv("RABBITMQ_PASSWORD", "guest")
	t.Setenv("NEW_DIR", dirs.newDir)
	t.Setenv("PROCESS_DIR", dirs.processDir)
	t.Setenv("DONE_DIR", dirs.doneDir)
	t.Setenv("ERROR_DIR", dirs.errorDir)
	t.Setenv("ERROR_LOG", filepath.Join(dirs.errorDir, "error.log"))
	t.Setenv("MOVIE_SOURCE", filepath.Join(root, "movies"))
	t.Setenv("SERIES_SOURCE", filepath.Join(root, "series"))

	return dirs
}
