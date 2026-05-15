package storage

import (
	"fmt"
	"os"

	"worker/internal/job"
)

func ValidateJobItems(
	job *job.MediaJob,
) error {

	for i := range job.Items {

		item := &job.Items[i]

		fullPath := job.ItemPath(item)

		info, err := os.Stat(fullPath)
		if err != nil {

			if os.IsNotExist(err) {
				return fmt.Errorf(
					"video item does not exist: %s",
					fullPath,
				)
			}

			return err
		}

		if info.IsDir() {
			return fmt.Errorf(
				"video item is directory: %s",
				fullPath,
			)
		}
	}

	return nil
}