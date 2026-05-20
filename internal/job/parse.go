package job

import (
	"encoding/json"
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

	for index, path := range msg.Videos {
		items = append(items, VideoItem{
			Index: index,

			FileName: filepath.Base(path),

		})
	}

	job := &MediaJob{
		ID:      msg.JobID,
		MediaID: msg.MediaID,
		Type:    msg.Type,
		Items:   items,
		Path:    jobPath(msg),
		State:   StateReceived,
	}

	return job, nil
}

func jobPath(msg Message) string {
	return filepath.Join(
		config.Load().NewDir,
		"job_"+msg.JobID,
	)
}
