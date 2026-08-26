# AGENTS.md

BudgetFlix worker: a single Go binary that consumes media upload jobs from RabbitMQ, converts videos to HLS via FFmpeg, and moves job folders through a filesystem state machine.

## Commands

- Run: `go run .` — requires a reachable RabbitMQ and `ffmpeg` on `PATH`.
- Build binary: `go build -o worker .`
- Test: `go test ./...` — no external services needed (uses `t.TempDir()`).
- Vet: `go vet ./...`
- No Makefile, linter config, or formatter config exists. There is no `gofmt` check in CI.

## Key structure

- `main.go` is the single entrypoint: load config → `storage.RecoverProcessingJobs` → connect RabbitMQ → `rabbitmq.Loop` in a goroutine → block on shutdown signals.
- Go module is named `worker` (NOT a domain path). All internal imports use the `worker/internal/...` prefix — `go get`/rename tools and external references must preserve this.
- `internal/rabbitmq` — connection, consumer, the `Loop` ack/nack loop. QoS prefetch is 1.
- `internal/handler.Media` — per-message handler, guarded by the package-level `PipelineWG` (`sync.WaitGroup`) used for graceful shutdown.
- `internal/pipeline` — movie processing. `Movie.Execute` runs inside `protected()` (WG-tracked).
- `internal/storage` — dir moves + file-state suffix changes (`.ready` → `.processing` → `.done` / `.error`).
- `internal/ffmpeg` — HLS command builder + runner. Runtime needs FFmpeg; tests only cover `cleanHLSOutput` (no binary).
- `internal/dto` — status messages published to the retry queue.

## Gotchas

- `config.Load()` is called fresh from multiple places (load paths, movie pipeline, main). Env is re-read each call; `Config.DataBase` is declared but never populated from env and never used — dead field, don't rely on it.
- Job folder resolution (`internal/job/parse.go`): looks at `NEW_DIR/job_<id>`, then falls back to `ERROR_DIR/job_<id>`. The message `path` field exists in the JSON shape but is not used; `NEW_DIR`/`ERROR_DIR` + `jobID` drive the location.
- Retry handling: the message has a `retry` bool. When `retry=true`, video paths resolve to `<file>.error` (in `ERROR_DIR`); otherwise to the original filename. Retry is started by the API layer — the worker never auto-retries.
- Only `MOVIE` is implemented; `SHOW` is defined but returns an error ("unsupported media type"). A `MOVIE` job must contain exactly one video item.
- On startup, `RecoverProcessingJobs` moves anything left in `PROCESS_DIR` to `ERROR_DIR`, marking every file `.error`.
- Success publishes `DONE` then `Ack`. Failure publishes `ERROR` then `Nack` with `requeue=false`. Both go to `VIDEO_UPLOAD_RETRY_QUEUE`.
- Shutdown waits up to 3 minutes (hardcoded `time.After(3 * time.Minute)` in `main.go`) for in-flight pipelines after `SIGTERM`/interrupt.

## Config / Docker

- Config is env-driven (`internal/config/config.go`); defaults are README-documented. `validate()` only enforces a subset — `RABBITMQ_USERNAME`, `RABBITMQ_PASSWORD`, `ERROR_LOG`, `MOVIE_SOURCE`, `SERIES_SOURCE` are not required-checked.
- MQ URL is built as `amqp://user:pass@host:port/`.
- Dockerfile: multi-stage, static build (`CGO_ENABLED=0`), runtime image `kokpeter/budgetflix-ffmpeg:latest`. FFmpeg must come from that image or be present in the container.
- CI (`.github/workflows/ci.yml`): runs on push to `main` only, builds and pushes `ghcr.io/budgetflix/budgetflix-worker:dev`. It does NOT run tests.
