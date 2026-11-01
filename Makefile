APP_NAME := Pegasus
WAILS3_VERSION := v3.0.0
GOPATH := $(shell go env GOPATH)
GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(GOPATH)/bin
endif
WAILS3 := $(shell command -v wails3 2>/dev/null || test ! -x "$(GOBIN)/wails3" || echo "$(GOBIN)/wails3")
HOST_ARCH := $(shell go env GOARCH)

.DEFAULT_GOAL := help
.PHONY: help deps doctor run dev build build-macos build-windows build-linux clean

help: ## Show available commands
	@printf "Pegasus build commands\n\n"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

deps: ## Install Wails v3 and download Go modules
	go install github.com/wailsapp/wails/v3/cmd/wails3@$(WAILS3_VERSION)
	go mod download

doctor: ## Check the Wails v3 build environment
	PATH="$(GOBIN):$$PATH" $(WAILS3) doctor

run: ## Run Pegasus from source
	go run .
dev: run

build: ## Build for the current OS and architecture
	PATH="$(GOBIN):$$PATH" $(WAILS3) build

build-macos: ## Build for macOS on the current architecture
	PATH="$(GOBIN):$$PATH" $(WAILS3) build GOOS=darwin GOARCH=$(HOST_ARCH)

build-windows: ## Build a Windows GUI executable (no console window)
	PATH="$(GOBIN):$$PATH" $(WAILS3) task windows:build ARCH=amd64

build-linux: ## Build for Linux on the current architecture
	PATH="$(GOBIN):$$PATH" $(WAILS3) build GOOS=linux GOARCH=$(HOST_ARCH)

clean: ## Remove generated binaries and distribution files
	rm -rf bin dist
	mkdir -p bin
	touch bin/.gitkeep
