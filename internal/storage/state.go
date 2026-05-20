package storage

import (
	"fmt"
	"os"
	"strings"
	"worker/internal/job"
)

type FileState string

const (
	StateNew        FileState = ""
	StateReady      FileState = ".ready"
	StateProcessing FileState = ".processing"
	StateDone       FileState = ".done"
	StateError      FileState = ".error"
)

func ChangeState(
	item *job.VideoItem,
	from FileState,
	to FileState,
) error {

	if from != StateNew &&
	!strings.HasSuffix(item.Path, string(from)) {
		return fmt.Errorf(
			"file does not have expected state suffix %s: %s",
			from,
			item.Path,
		)
	}

	newPath := strings.TrimSuffix(
		item.Path,
		string(from),
	) + string(to)

	if err := os.Rename(item.Path, newPath); err != nil {
		return err
	}

	item.Path = newPath

	item.FileName = strings.TrimSuffix(
		item.FileName,
		string(from),
	) + string(to)

	return nil
}