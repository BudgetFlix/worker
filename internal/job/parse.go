package job

import (
	"encoding/json"
	"os"
	"path/filepath"

	"worker/internal/config"
)

type MediaType string

const (
	MediaTypeMovie MediaType = "MOVIE"
	MediaTypeShow  MediaType = "SHOW"
)

type Message struct {
	JobID   string         `json:"jobID"`
	MediaID int64          `json:"mediaID"`
	Videos  map[int]string `json:"videos"`
	Type    MediaType      `json:"type"`
	Path    string         `json:"path"`
}

func FromJSON(data []byte) (*MediaJob, error) {
	var msg Message

	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	items := make([]VideoItem, 0, len(msg.Videos))

	jobPath := jobPathMaker(msg)

	for index, path := range msg.Videos {
		filename := filepath.Base(path)
		items = append(items, VideoItem{
			Index:    index,
			FileName: filename,
			Path:     filepath.Join(jobPath, filename),
		})
	}

	job := &MediaJob{
		ID:      msg.JobID,
		MediaID: msg.MediaID,
		Type:    msg.Type,
		Items:   items,
		Path:    jobPath,
		State:   StateReceived,
	}

	return job, nil
}

func jobPathMaker(msg Message) string {
	cfg := config.Load()
	newJobPath := filepath.Join(
		cfg.NewDir,
		"job_"+msg.JobID,
	)

	if pathExists(newJobPath) {
		return newJobPath
	}

	errorJobPath := filepath.Join(
		cfg.ErrorDir,
		"job_"+msg.JobID,
	)

	// TODO: Resolve the job folder from the incoming message once the API sends it.
	if pathExists(errorJobPath) {
		return errorJobPath
	}

	return newJobPath
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
