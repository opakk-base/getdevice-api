# Docker build environment for cross-compiling Pegasus to Linux.
#
# Wails v2 cannot natively cross-compile to Linux from macOS/Windows because
# the binary links GTK/WebKit at build time. This image provides those dev
# libraries plus the Wails CLI, so `make build-linux-docker` works anywhere
# Docker runs.
#
# Build (from repo root):
#   docker build -f build/linux/Dockerfile -t getdevice-linux-builder .
#
# Usage is wired up in the Makefile (`make build-linux-docker`).

FROM golang:1.26-trixie

# Build dependencies required by Wails v3 on Linux.
RUN apt-get update && apt-get install -y --no-install-recommends \
        pkg-config \
        gcc \
        libc6-dev \
        libgtk-4-dev \
        libwebkitgtk-6.0-dev \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

WORKDIR /project
