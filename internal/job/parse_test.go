package job

import (
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
