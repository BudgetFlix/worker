# =========================================================
# GO BUILD
# =========================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o worker .

# =========================================================
# FINAL IMAGE
# =========================================================
FROM kokpeter/budgetflix-ffmpeg:latest

WORKDIR /app

COPY --from=builder /app/worker .

ENTRYPOINT ["./worker"]