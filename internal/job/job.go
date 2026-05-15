package job

import "path/filepath"
import "time"

type MediaJob struct {
	ID string

	MediaID int64

	Type MediaType

	Items []VideoItem

	Path string

	State State

	CreatedAt time.Time
	UpdatedAt time.Time
}


type VideoItem struct {
	Index int

	FileName string

	State State

	Error string
}



func (j *MediaJob) ItemPath(
	item *VideoItem,
) string {
	return filepath.Join(
		j.Path,
		item.FileName,
	)
}