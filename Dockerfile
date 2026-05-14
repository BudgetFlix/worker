# ===== Build stage =====
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Go modulok
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Binary build
RUN CGO_ENABLED=0 GOOS=linux go build -o worker .

# ===== Runtime stage =====
FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

# Runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Binary copy
COPY --from=builder /app/worker worker


# Execute
ENTRYPOINT ["./worker"]