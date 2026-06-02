package ffmpeg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanHLSOutputRemovesPreviousPlaylistAndSegments(t *testing.T) {
	outputDir := t.TempDir()

	for _, file := range []string{
		"index.m3u8",
		"seg_000.ts",
		"seg_001.ts",
		"cover.jpg",
	} {
		if err := os.WriteFile(filepath.Join(outputDir, file), []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := cleanHLSOutput(outputDir); err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"index.m3u8",
		"seg_000.ts",
		"seg_001.ts",
	} {
		if _, err := os.Stat(filepath.Join(outputDir, file)); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, got error: %v", file, err)
		}
	}

	if _, err := os.Stat(filepath.Join(outputDir, "cover.jpg")); err != nil {
		t.Fatal(err)
	}
}
