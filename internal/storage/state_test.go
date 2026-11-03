package storage

import (
	"os"
	"path/filepath"
	"testing"

	"worker/internal/job"
)

func TestChangeToErrorIgnoresAlreadyErroredFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "movie.mp4.error")

	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	item := &job.VideoItem{
		FileName: "movie.mp4.error",
		Path:     path,
	}

	if err := ChangeToError(item); err != nil {
		t.Fatal(err)
	}

	if item.Path != path {
		t.Fatalf("expected path to stay unchanged, got %s", item.Path)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestRetryStateTransitionRemovesErrorSuffix(t *testing.T) {
	dir := t.TempDir()
	errorPath := filepath.Join(dir, "movie.mp4.error")

	if err := os.WriteFile(errorPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	item := &job.VideoItem{
		FileName: "movie.mp4.error",
		Path:     errorPath,
	}

	if err := ChangeState(item, StateError, StateReady); err != nil {
		t.Fatal(err)
	}

	expectedPath := filepath.Join(dir, "movie.mp4.ready")
	if item.Path != expectedPath {
		t.Fatalf("expected path %s, got %s", expectedPath, item.Path)
	}
	if item.FileName != "movie.mp4.ready" {
		t.Fatalf("expected filename movie.mp4.ready, got %s", item.FileName)
	}
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatal(err)
	}
}
