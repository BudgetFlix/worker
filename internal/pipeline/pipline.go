package pipeline

import (
	"context"

	"worker/internal/job"
)

type Pipeline interface {
	Execute(
		ctx context.Context,
		job *job.MediaJob,
	) error
}