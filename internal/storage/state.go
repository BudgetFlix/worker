package storage

import (
	"fmt"
	"os"
	"strings"
)

type FileState string

const (
	StateNew        FileState = ".new"
	StateReady      FileState = ".ready"
	StateProcessing FileState = ".processing"
	StateDone       FileState = ".done"
	StateError      FileState = ".error"
)

func ChangeState(
	path string,
	from FileState,
	to FileState,
) (string, error) {

	if !strings.HasSuffix(path, string(from)) {
		return "", fmt.Errorf(
			"file does not have expected state suffix %s: %s",
			from,
			path,
		)
	}

	newPath := strings.TrimSuffix(path, string(from)) + string(to)

	err := os.Rename(path, newPath)
	if err != nil {
		return "", err
	}

	return newPath, nil
}