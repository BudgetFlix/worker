package handler

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"

	"worker/internal/job"
)

func TestFailMediaJobMovesNewJobToErrorDir(t *testing.T) {
	root := t.TempDir()
	newDir := filepath.Join(root, "new")
	errorDir := filepath.Join(root, "error")
	jobPath := filepath.Join(newDir, "job_123")
	videoPath := filepath.Join(jobPath, "movie.mp4")

	if err := os.MkdirAll(jobPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(videoPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ERROR_DIR", errorDir)

	mediaJob := &job.MediaJob{
		ID:   "123",
		Path: jobPath,
		Items: []job.VideoItem{{
			FileName: "movie.mp4",
			Path:     videoPath,
		}},
	}
	cause := errors.New("unsupported media type")

	if err := failMediaJob(mediaJob, cause); !errors.Is(err, cause) {
		t.Fatalf("expected original error, got %v", err)
	}

	errorFilePath := filepath.Join(errorDir, "job_123", "movie.mp4.error")
	if _, err := os.Stat(errorFilePath); err != nil {
		t.Fatalf("expected failed file at %s: %v", errorFilePath, err)
	}
	if _, err := os.Stat(jobPath); !os.IsNotExist(err) {
		t.Fatalf("expected new job to be moved, got %v", err)
	}
}

func TestFailMediaJobDoesNotMoveRetryJobOutOfErrorDir(t *testing.T) {
	root := t.TempDir()
	errorDir := filepath.Join(root, "error")
	jobPath := filepath.Join(errorDir, "job_123")
	errorPath := filepath.Join(jobPath, "movie.mp4.error")

	if err := os.MkdirAll(jobPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(errorPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ERROR_DIR", errorDir)

	mediaJob := &job.MediaJob{
		ID:   "123",
		Path: jobPath,
		Items: []job.VideoItem{{
			FileName: "movie.mp4.error",
			Path:     errorPath,
		}},
	}

	if err := failMediaJob(mediaJob, errors.New("retry failed")); err == nil {
		t.Fatal("expected retry failure error")
	}
	if _, err := os.Stat(errorPath); err != nil {
		t.Fatalf("expected retry file to remain in error directory: %v", err)
	}
}

func TestMediaAcceptsCompletedJobRedelivery(t *testing.T) {
	root := t.TempDir()
	doneDir := filepath.Join(root, "done")
	if err := os.MkdirAll(filepath.Join(doneDir, "job_123"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DONE_DIR", doneDir)
	t.Setenv("NEW_DIR", filepath.Join(root, "new"))
	t.Setenv("ERROR_DIR", filepath.Join(root, "error"))

	msg := amqp.Delivery{Body: []byte(`{"jobID":"123","mediaID":1,"type":"MOVIE","videos":{"0":"movie.mp4"}}`)}
	if err := Media(msg); err != nil {
		t.Fatalf("expected completed redelivery to be accepted: %v", err)
	}
}
