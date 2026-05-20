# =========================================================
# FFMPEG + X264 BUILD
# =========================================================
FROM alpine:3.22 AS ffmpeg-builder

RUN apk add --no-cache \
    bash \
    git \
    build-base \
    yasm \
    nasm \
    pkgconfig \
    wget \
    tar \
    xz

WORKDIR /build

# =========================================================
# BUILD X264
# =========================================================
RUN git clone --depth 1 https://code.videolan.org/videolan/x264.git

WORKDIR /build/x264

RUN ./configure \
    --prefix="/opt/ffmpeg" \
    --enable-static \
    --disable-opencl && \
    make -j$(nproc) && \
    make install

# =========================================================
# BUILD FFMPEG
# =========================================================
WORKDIR /build

RUN wget https://ffmpeg.org/releases/ffmpeg-7.1.tar.xz && \
    tar -xJf ffmpeg-7.1.tar.xz

WORKDIR /build/ffmpeg-7.1

RUN PKG_CONFIG_PATH="/opt/ffmpeg/lib/pkgconfig" \
    ./configure \
    \
    --prefix="/opt/ffmpeg" \
    \
    --pkg-config-flags="--static" \
    --extra-cflags="-I/opt/ffmpeg/include" \
    --extra-ldflags="-L/opt/ffmpeg/lib" \
    \
    --enable-gpl \
    --enable-static \
    --disable-shared \
    \
    --disable-everything \
    \
    --enable-ffmpeg \
    \
    --enable-protocol=file \
    --enable-protocol=pipe \
    \
    --enable-demuxer=mov \
    --enable-demuxer=matroska \
    --enable-demuxer=avi \
    \
    --enable-muxer=hls \
    --enable-muxer=mpegts \
    \
    --enable-parser=h264 \
    --enable-parser=aac \
    \
    --enable-decoder=h264 \
    --enable-decoder=aac \
    \
    --enable-encoder=aac \
    \
    --enable-libx264 \
    --enable-encoder=libx264 \
    \
    --enable-filter=scale \
    --enable-swscale \
    \
    --enable-small \
    \
    --disable-doc \
    --disable-debug \
    --disable-network \
    --disable-avdevice \
    --disable-postproc && \
    make -j$(nproc) && \
    make install

# =========================================================
# BUILD GO WORKER
# =========================================================
FROM golang:1.24-alpine AS go-builder

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
FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=ffmpeg-builder /opt/ffmpeg /opt/ffmpeg
COPY --from=go-builder /app/worker .

ENV PATH="/opt/ffmpeg/bin:${PATH}"

ENTRYPOINT ["./worker"]