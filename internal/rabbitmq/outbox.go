package rabbitmq

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"worker/internal/config"
	"worker/internal/dto"
)

type statusOutboxEntry struct {
	Queue   string            `json:"queue"`
	Message dto.MediaRetryMsg `json:"message"`
}

func PersistStatus(
	cfg *config.Config,
	queue string,
	message dto.MediaRetryMsg,
) error {
	dir := statusOutboxDir(cfg)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	body, err := json.Marshal(statusOutboxEntry{Queue: queue, Message: message})
	if err != nil {
		return err
	}

	path := statusOutboxPath(cfg, message)
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".status-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

func RemoveStatus(
	cfg *config.Config,
	message dto.MediaRetryMsg,
) error {
	err := os.Remove(statusOutboxPath(cfg, message))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func FlushStatusOutbox(
	ctx context.Context,
	cfg *config.Config,
	producer *Producer,
) error {
	entries, err := os.ReadDir(statusOutboxDir(cfg))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var errs []error
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(statusOutboxDir(cfg), entry.Name())
		body, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("read status outbox %s: %w", path, err))
			continue
		}

		var outboxEntry statusOutboxEntry
		if err := json.Unmarshal(body, &outboxEntry); err != nil {
			errs = append(errs, fmt.Errorf("decode status outbox %s: %w", path, err))
			continue
		}

		select {
		case <-ctx.Done():
			return errors.Join(append(errs, ctx.Err())...)
		default:
		}

		if err := producer.Publish(outboxEntry.Queue, outboxEntry.Message); err != nil {
			errs = append(errs, fmt.Errorf("publish status outbox %s: %w", path, err))
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("remove status outbox %s: %w", path, err))
		}
	}

	return errors.Join(errs...)
}

func RetryStatusOutbox(
	ctx context.Context,
	cfg *config.Config,
	producer *Producer,
) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := FlushStatusOutbox(ctx, cfg, producer); err != nil {
				fmt.Printf("failed to flush status outbox: %v\n", err)
			}
		}
	}
}

func statusOutboxDir(cfg *config.Config) string {
	return filepath.Join(cfg.ErrorDir, ".status-outbox")
}

func statusOutboxPath(cfg *config.Config, message dto.MediaRetryMsg) string {
	hash := sha256.Sum256([]byte(message.ID + "\x00" + message.Status))
	return filepath.Join(statusOutboxDir(cfg), hex.EncodeToString(hash[:])+".json")
}
