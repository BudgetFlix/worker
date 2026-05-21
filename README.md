<div align=center>

  

# BudgetFlix Worker

**Lightweight Go worker service** for the BudgetFlix backend. It consumes media upload jobs from **RabbitMQ**, validates the incoming job folder, converts movie files to **HLS with FFmpeg**, and moves job folders through a **filesystem-based lifecycle**.

<br>

## Tech Stack

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![FFmpeg](https://img.shields.io/badge/FFmpeg-007808?style=for-the-badge&logo=ffmpeg&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)

</div>

## What It Does

- Consumes media jobs from **RabbitMQ**
- Validates incoming **job payloads**
- Processes `MOVIE` jobs
- Converts videos into **HLS streams with FFmpeg**
- Tracks **file processing states**
- Moves jobs through **lifecycle directories**
- Handles RabbitMQ **ACK/NACK** flow
- Supports **graceful shutdown**



## Project Structure

```text
.
+-- main.go                    # Worker entrypoint and graceful shutdown
+-- Dockerfile                 # Multi-stage build with FFmpeg runtime image
+-- internal
|   +-- config                 # Environment configuration
|   +-- ffmpeg                 # HLS command builder and runner
|   +-- handler                # RabbitMQ message handler
|   +-- job                    # Job parsing, validation and models
|   +-- logger                 # Logging helpers
|   +-- pipeline               # Media processing pipelines
|   +-- rabbitmq               # Connection, consumer, ack/nack loop
|   +-- storage                # File state and directory movement helpers
+-- README.md
```

## Worker Flow

```mermaid
flowchart LR
    RabbitMQ["RabbitMQ: video.upload.queue"]
    Handler["Media handler"]
    Validate["Parse and validate job"]
    Ready["Mark item as .ready"]
    ProcessDir["Move job to PROCESS_DIR"]
    Processing["Mark item as .processing"]
    HLS["FFmpeg HLS encoding"]
    Done["Mark item as .done"]
    DoneDir["Move job to DONE_DIR"]
    Error["Mark item as .error and move to ERROR_DIR"]
    Nack["Nack failed message"]

    RabbitMQ --> Handler --> Validate --> Ready --> ProcessDir --> Processing --> HLS
    HLS --> Done --> DoneDir
    Handler --> Nack
    Processing --> Error
    HLS --> Error
    Error --> Nack
```

## Job Message Format

The worker expects JSON messages with this shape:

```json
{
  "jobID": "example-job-id",
  "mediaID": 123,
  "type": "MOVIE",
  "path": "optional/source/path/from-producer",
  "videos": {
    "0": "/path/to/video.mp4"
  }
}
```

Notes:

- `type` currently supports `MOVIE` in the implemented pipeline.
- `MOVIE` jobs must contain exactly one video item.
- The worker derives the actual job directory from `NEW_DIR/job_<jobID>`.
- Video file names are taken from the submitted video paths and resolved inside the job directory.

## File Lifecycle

For a job with ID `abc123`, the worker expects the upload folder to exist at:

```text
NEW_DIR/job_abc123
```

During processing, the job directory and file states move through:

```text
NEW_DIR/job_abc123
PROCESS_DIR/job_abc123
DONE_DIR/job_abc123
```

On failure, the job is moved to:

```text
ERROR_DIR/job_abc123
```

The video file itself is renamed as it advances:

```text
movie.mp4
movie.mp4.ready
movie.mp4.processing
movie.mp4.done
```

If processing fails while the item is in progress, it is marked as:

```text
movie.mp4.error
```

## HLS Output

For movie jobs, output is created at:

```text
MOVIE_SOURCE/<mediaID>/hls
```

The FFmpeg command uses:

- H.264 video (`libx264`)
- AAC audio
- HLS VOD playlist
- 6 second HLS segments
- `index.m3u8` as the playlist file
- `seg_%03d.ts` as the segment pattern

## Configuration

Configuration is read from environment variables. Missing variables fall back to local defaults.

| Variable | Default | Description |
| --- | --- | --- |
| `RABBITMQ_HOST` | `localhost` | RabbitMQ host |
| `RABBITMQ_PORT` | `5672` | RabbitMQ port |
| `RABBITMQ_USERNAME` | `guest` | RabbitMQ username |
| `RABBITMQ_PASSWORD` | `guest` | RabbitMQ password |
| `NEW_DIR` | `/tmp/new` | Directory where new job folders are expected |
| `PROCESS_DIR` | `/tmp/process` | Directory for active jobs |
| `DONE_DIR` | `/tmp/done` | Directory for completed jobs |
| `ERROR_DIR` | `/tmp/error` | Directory for failed jobs |
| `ERROR_LOG` | `/tmp/error/error.log` | Error log path placeholder |
| `MOVIE_SOURCE` | `/tmp/movies` | Movie library root |
| `SERIES_SOURCE` | `/tmp/series` | Series library root |

The RabbitMQ connection URL is built as:

```text
amqp://<RABBITMQ_USERNAME>:<RABBITMQ_PASSWORD>@<RABBITMQ_HOST>:<RABBITMQ_PORT>/
```

## Running Locally

Prerequisites:

- Go `1.24.3`
- RabbitMQ
- FFmpeg available on `PATH`

Run the worker:

```bash
go run .
```

Run tests:

```bash
go test ./...
```

Example local environment:

```bash
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USERNAME=guest
export RABBITMQ_PASSWORD=guest
export NEW_DIR=/tmp/new
export PROCESS_DIR=/tmp/process
export DONE_DIR=/tmp/done
export ERROR_DIR=/tmp/error
export MOVIE_SOURCE=/tmp/movies
```

## Docker

Build the image:

```bash
docker build -t budgetflix-worker .
```

Run the container:

```bash
docker run --rm \
  -e RABBITMQ_HOST=host.docker.internal \
  -e RABBITMQ_PORT=5672 \
  -e RABBITMQ_USERNAME=guest \
  -e RABBITMQ_PASSWORD=guest \
  -v /tmp/new:/tmp/new \
  -v /tmp/process:/tmp/process \
  -v /tmp/done:/tmp/done \
  -v /tmp/error:/tmp/error \
  -v /tmp/movies:/tmp/movies \
  budgetflix-worker
```

The Dockerfile builds a static Go binary in a Go Alpine builder image and runs it in an FFmpeg-enabled runtime image.

## RabbitMQ Behavior

- Queue name: `video.upload.queue`
- Manual acknowledgements are used.
- QoS prefetch is set to `1`, so the worker processes one message at a time per consumer.
- Successful jobs are acknowledged with `Ack`.
- Failed jobs are rejected with `Nack` and `requeue=false`.

## Graceful Shutdown

The worker listens for interrupt and `SIGTERM` signals. On shutdown it:

1. Stops the consumer loop.
2. Waits for active pipelines to finish.
3. Exits after all pipelines complete or after a 3 minute timeout.

## Current Limitations

- Only the `MOVIE` pipeline is implemented.
- `SHOW` is defined as a media type but is not processed yet.
- Movie jobs must contain exactly one video file.
- Retry handling, metrics and health checks are not implemented yet.
- The worker assumes the producer has already created the expected job folder and source file under `NEW_DIR`.

## Development Notes

Useful commands:

```bash
go test ./...
go run .
go build -o worker .
```

Before changing processing behavior, check these areas:

- `internal/job` for message format and validation
- `internal/storage` for job movement and file state changes
- `internal/pipeline` for media-specific processing
- `internal/ffmpeg` for HLS command generation
- `internal/rabbitmq` for queue consumption behavior
