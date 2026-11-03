package job

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVideoPathMakerRetryUsesErrorStateFile(t *testing.T) {
	msg := Message{Retry: true}

	path := VideoPathMaker(msg, "movie.mp4", "/tmp/error/job_123")

	expected := filepath.Join("/tmp/error/job_123", "movie.mp4.error")
	if path != expected {
		t.Fatalf("expected %s, got %s", expected, path)
	}
}

func TestVideoPathMakerUsesErrorStateForRedelivery(t *testing.T) {
	dir := t.TempDir()
	errorPath := filepath.Join(dir, "movie.mp4.error")
	if err := os.WriteFile(errorPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	path := VideoPathMaker(Message{}, "movie.mp4", dir)
	if path != errorPath {
		t.Fatalf("expected %s, got %s", errorPath, path)
	}
}

func TestFromJSONResolvesCompletedJobForRedelivery(t *testing.T) {
	root := t.TempDir()
	doneDir := filepath.Join(root, "done")
	jobPath := filepath.Join(doneDir, "job_123")
	if err := os.MkdirAll(jobPath, 0755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("DONE_DIR", doneDir)
	t.Setenv("NEW_DIR", filepath.Join(root, "new"))
	t.Setenv("ERROR_DIR", filepath.Join(root, "error"))

	job, err := FromJSON([]byte(`{"jobID":"123","mediaID":1,"type":"MOVIE","videos":{"0":"movie.mp4"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if job.Path != jobPath {
		t.Fatalf("expected completed job path %s, got %s", jobPath, job.Path)
	}
}
