package rabbitmq

import (
	"os"
	"path/filepath"
	"testing"

	"worker/internal/config"
	"worker/internal/dto"
)

func TestPersistStatusCreatesDurableEntry(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{ErrorDir: filepath.Join(root, "error")}
	message := dto.MediaRetryMsg{ID: "123", Status: "DONE"}

	if err := PersistStatus(cfg, "status.queue", message); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(filepath.Join(cfg.ErrorDir, ".status-outbox"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one outbox entry, got %d", len(entries))
	}
}

func TestRemoveStatusRemovesDurableEntry(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{ErrorDir: filepath.Join(root, "error")}
	message := dto.MediaRetryMsg{ID: "123", Status: "ERROR"}

	if err := PersistStatus(cfg, "status.queue", message); err != nil {
		t.Fatal(err)
	}
	if err := RemoveStatus(cfg, message); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(statusOutboxPath(cfg, message)); !os.IsNotExist(err) {
		t.Fatalf("expected outbox entry to be removed, got %v", err)
	}
}
