package storage

import (
	"fmt"
	"os"
)

func ValidateJobDirectory(
	path string,
) error {

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"job directory does not exist: %s",
				path,
			)
		}

		return err
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"job path is not directory: %s",
			path,
		)
	}

	return nil
}