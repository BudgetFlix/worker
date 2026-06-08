package job

import (
	"errors"

	"os"
)

func Validate(job *MediaJob) error {

	if job.ID == "" {
		return errors.New("missing job id")
	}

	if job.MediaID <= 0 {
		return errors.New("invalid media id")
	}

	if job.Path == "" {
		return errors.New("missing job path")
	}

	if len(job.Items) == 0 {
		return errors.New("job has no video items")
	}

	info, err := os.Stat(job.Path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("job path is not directory")
	}

	for i := range job.Items {
    item := &job.Items[i]

    info, err := os.Stat(item.Path)
    if err != nil {
        return err
    }

    if info.IsDir() {
        return errors.New("video item is directory")
    }
}

	return nil
}
