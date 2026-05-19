package job

import (
	"path/filepath"
	"testing"
)

func TestFromJSONBuildsJobPathAndKeepsFileName(t *testing.T) {
	t.Setenv("RABBITMQ_HOST", "localhost")
	t.Setenv("RABBITMQ_PORT", "5672")
	t.Setenv("RABBITMQ_USERNAME", "guest")
	t.Setenv("RABBITMQ_PASSWORD", "guest")
	t.Setenv("NEW_DIR", "/media/new")
	t.Setenv("PROCESS_DIR", "/media/process")
	t.Setenv("DONE_DIR", "/media/done")
	t.Setenv("ERROR_DIR", "/media/error")
	t.Setenv("ERROR_LOG", "/media/error/error.log")
	t.Setenv("MOVIE_SOURCE", "/media/movies")
	t.Setenv("SERIES_SOURCE", "/media/series")

	job, err := FromJSON([]byte(`{
		"jobID": "1234",
		"mediaID": 42,
		"type": "MOVIE",
		"videos": {
			"0": "/uploads/movie.mp4"
		}
	}`))
	if err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}

	expectedPath := filepath.Join("/media/new", "job_1234")
	if job.Path != expectedPath {
		t.Fatalf("expected job path %q, got %q", expectedPath, job.Path)
	}

	if len(job.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(job.Items))
	}

	if job.Items[0].FileName != "movie.mp4" {
		t.Fatalf("expected file name %q, got %q", "movie.mp4", job.Items[0].FileName)
	}
}

func TestFromJSONKeepsReadyFileName(t *testing.T) {
	t.Setenv("NEW_DIR", "/media/new")

	job, err := FromJSON([]byte(`{
		"jobID": "1234",
		"mediaID": 42,
		"type": "MOVIE",
		"videos": {
			"0": "movie.mp4.ready"
		},
		"path": "/custom/job_1234"
	}`))
	if err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}

	expectedPath := filepath.Join("/media/new", "job_1234")
	if job.Path != expectedPath {
		t.Fatalf("expected job path %q, got %q", expectedPath, job.Path)
	}

	if job.Items[0].FileName != "movie.mp4.ready" {
		t.Fatalf("expected ready file name to be preserved, got %q", job.Items[0].FileName)
	}
}
