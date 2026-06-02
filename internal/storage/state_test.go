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
