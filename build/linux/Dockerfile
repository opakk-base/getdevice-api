# Docker build environment for cross-compiling GetDevice to Linux.
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

FROM golang:1.25-trixie

# Build dependencies required by Wails v2 on Linux.
# Node.js comes from the Debian repo; it is unused when frontend commands are
# empty in wails.json, but including it keeps `wails doctor` happy.
RUN apt-get update && apt-get install -y --no-install-recommends \
        nodejs \
        npm \
        pkg-config \
        gcc \
        libc6-dev \
        libgtk-3-dev \
        libwebkit2gtk-4.1-dev \
        libayatana-appindicator3-dev \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest

WORKDIR /project
