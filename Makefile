# =============================================================================
#  GetDevice — Makefile
#  Builds, packages, and installs the Wails (Go + JS) desktop application on
#  macOS, Windows, and Linux.
#
#  Usage:
#    make help              Show this help
#    make run               Run in development mode (hot reload)
#    make build             Build for the current OS/arch
#    make build-macos       Build a macOS universal binary (.app)
#    make build-windows     Build a Windows amd64 binary (.exe)
#    make build-linux       Build a Linux amd64 binary
#    make build-all         Build all three platforms (native + cross)
#
#  Installers:
#    make installer-macos   Build .app and package a DMG (macOS host)
#    make installer-windows Build NSIS setup .exe (Windows host, needs makensis)
#    make installer-linux   Build .deb / .rpm / .tar.gz (needs nfpm)
#    make installer         Package for the detected host OS
#
#    make clean             Remove build artifacts
# =============================================================================

# --- Project metadata --------------------------------------------------------
APP_NAME      := GetDevice
BUNDLE_NAME   := getdevice-api.app
OUTPUT_DIR    := build/bin
DIST_DIR      := dist
VERSION       ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
# Plain numeric version (strip leading v and git suffixes) — required by nfpm
PKG_VERSION   := $(shell echo $(VERSION) | sed -E 's/^v//; s/-.*//')
ifeq ($(PKG_VERSION),)
  PKG_VERSION := 0.0.0
endif

# --- OS detection ------------------------------------------------------------
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows_NT)

ifeq ($(OS),Windows_NT)
  OS_DETECTED := Windows
else
  ifneq ($(findstring MINGW,$(UNAME_S)),)
    OS_DETECTED := Windows
  else ifneq ($(findstring MSYS,$(UNAME_S)),)
    OS_DETECTED := Windows
  else ifneq ($(findstring CYGWIN,$(UNAME_S)),)
    OS_DETECTED := Windows
  else ifneq ($(findstring Darwin,$(UNAME_S)),)
    OS_DETECTED := macOS
  else ifneq ($(findstring Linux,$(UNAME_S)),)
    OS_DETECTED := Linux
  else
    OS_DETECTED := Unknown
  endif
endif

# --- Architecture detection ---------------------------------------------------
UNAME_M := $(shell uname -m 2>/dev/null || echo unknown)
ifeq ($(UNAME_M),arm64)
  HOST_ARCH := arm64
else ifeq ($(UNAME_M),aarch64)
  HOST_ARCH := arm64
else ifeq ($(UNAME_M),x86_64)
  HOST_ARCH := amd64
else
  HOST_ARCH := $(UNAME_M)
endif

# --- Executable extension (Windows only) --------------------------------------
ifeq ($(OS_DETECTED),Windows)
  EXE_EXT := .exe
else
  EXE_EXT :=
endif

# --- Wails platform strings ---------------------------------------------------
WAILS_DARWIN  := darwin/universal
WAILS_WINDOWS := windows/amd64
WAILS_LINUX   := linux/amd64

# --- Wails CLI resolution -----------------------------------------------------
# `go install` places wails in $(go env GOPATH)/bin (or GOBIN), which is often
# not on PATH in a fresh shell. Resolve the binary explicitly so `make deps`
# followed by `make dev`/`make build` works without re-sourcing your profile.
GOPATH := $(shell go env GOPATH 2>/dev/null)
GOBIN  := $(shell go env GOBIN 2>/dev/null)

ifeq ($(GOBIN),)
  GOBIN := $(GOPATH)/bin
endif

WAILS := $(shell command -v wails 2>/dev/null)
ifeq ($(WAILS),)
  ifneq ($(wildcard $(GOBIN)/wails$(EXE_EXT)),)
    WAILS := $(GOBIN)/wails$(EXE_EXT)
  endif
endif

# =============================================================================
#  Default target
# =============================================================================
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@printf "\033[1mGetDevice Makefile\033[0m (detected: $(OS_DETECTED)/$(HOST_ARCH))\n\n"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@printf "\nVersion: $(VERSION)\n"

# =============================================================================
#  Setup & checks
# =============================================================================
.PHONY: deps
deps: ## Download Go modules and verify toolchain (wails, node)
	@command -v go >/dev/null 2>&1 || (echo "ERROR: Go is not installed (go 1.22+ required)" && exit 1)
	@command -v node >/dev/null 2>&1 || (echo "ERROR: Node.js is not installed (Node 15+ required)" && exit 1)
	@WAILS_BIN="$$(command -v wails 2>/dev/null || true)"; \
	if [ -z "$$WAILS_BIN" ] && [ -x "$(GOBIN)/wails$(EXE_EXT)" ]; then \
		WAILS_BIN="$(GOBIN)/wails$(EXE_EXT)"; \
	fi; \
	if [ -z "$$WAILS_BIN" ]; then \
		echo "Wails CLI not found. Installing to $(GOBIN)..."; \
		go install github.com/wailsapp/wails/v2/cmd/wails@latest; \
		WAILS_BIN="$(GOBIN)/wails$(EXE_EXT)"; \
	fi; \
	if [ ! -x "$$WAILS_BIN" ]; then \
		echo "ERROR: Wails install failed (expected binary at $$WAILS_BIN)."; \
		echo "Check that $(GOBIN) is writable and on your PATH."; \
		exit 1; \
	fi
	go mod download
	@echo "Dependencies ready. Wails CLI: $(GOBIN)/wails$(EXE_EXT)"

.PHONY: doctor
doctor: ## Run `wails doctor` to verify the build environment
	@if [ -z "$(WAILS)" ]; then \
		echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; \
	fi
	$(WAILS) doctor

.PHONY: check-env
check-env: ## Show the detected build environment
	@echo "OS          : $(OS_DETECTED)"
	@echo "Host arch   : $(HOST_ARCH)"
	@echo "Binary ext  : $(EXE_EXT)"
	@echo "Version     : $(VERSION)"
	@echo "Output dir  : $(OUTPUT_DIR)"

# =============================================================================
#  Development
# =============================================================================
.PHONY: run dev
run: ## Run the app in development mode (hot reload)
dev: run
run:
	@if [ -z "$(WAILS)" ]; then \
		echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; \
	fi
	$(WAILS) dev

# =============================================================================
#  Build
# =============================================================================
.PHONY: build
build: ## Build for the current OS/arch (native)
	@echo "Building for $(OS_DETECTED)/$(HOST_ARCH)..."
	@if [ -z "$(WAILS)" ]; then \
		echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; \
	fi
	$(WAILS) build -s -skipbindings

.PHONY: build-macos
build-macos: ## Build a macOS universal binary (.app)
	@echo "Building macOS universal binary..."
	@if [ -z "$(WAILS)" ]; then \
		echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; \
	fi
	@echo "NOTE: macOS builds must be performed on macOS (or via CI/remote Mac)."
	$(WAILS) build -platform $(WAILS_DARWIN) -s -skipbindings

.PHONY: build-windows
build-windows: ## Build a Windows amd64 binary (.exe)
	@echo "Building Windows amd64 binary..."
	@if [ -z "$(WAILS)" ]; then \
		echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; \
	fi
	@if [ "$(OS_DETECTED)" != "Windows" ]; then \
		echo "NOTE: Cross-compiling Windows from $(OS_DETECTED) may require extra setup"; \
		echo "      (Docker on Linux, or a Windows machine/CI)."; \
	fi
	$(WAILS) build -platform $(WAILS_WINDOWS) -s -skipbindings

.PHONY: build-linux
build-linux: ## Build a Linux amd64 binary
	@echo "Building Linux amd64 binary..."
	@if [ -z "$(WAILS)" ]; then \
		echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; \
	fi
	@if [ "$(OS_DETECTED)" != "Linux" ]; then \
		echo "NOTE: Cross-compiling Linux from $(OS_DETECTED) may require extra setup."; \
	fi
	$(WAILS) build -platform $(WAILS_LINUX) -s -skipbindings

.PHONY: build-all
build-all: ## Build for macOS, Windows, and Linux
	@echo "Building all platforms..."
	$(MAKE) build-macos
	$(MAKE) build-windows
	$(MAKE) build-linux
	@echo "All builds complete. Artifacts are in $(OUTPUT_DIR)/"

# --- Linux cross-build via Docker (Wails cannot cross-compile to Linux natively)
# Uses our builder image (build/linux/Dockerfile): Go + Node + GTK/WebKit dev libs
# + Wails CLI. Single-target linux builds output build/bin/GetDevice.
LINUX_BUILDER_IMAGE := getdevice-linux-builder

.PHONY: build-linux-docker
build-linux-docker: ## Build the Linux amd64 binary inside Docker (works from macOS/Windows)
	@command -v docker >/dev/null 2>&1 || { \
		echo "ERROR: docker is required for cross-building Linux from $(OS_DETECTED)."; \
		echo "       Install Docker Desktop, or build on a Linux machine/CI."; exit 1; \
	}
	docker info >/dev/null 2>&1 || { \
		echo "ERROR: docker daemon is not running. Start Docker Desktop first."; exit 1; \
	}
	docker image inspect $(LINUX_BUILDER_IMAGE) >/dev/null 2>&1 && \
		[ "$$(docker image inspect $(LINUX_BUILDER_IMAGE) --format '{{.Architecture}}')" = "amd64" ] || \
		docker build --platform linux/amd64 -f build/linux/Dockerfile -t $(LINUX_BUILDER_IMAGE) .
	docker run --rm --platform linux/amd64 \
		-v "$(CURDIR):/project" \
		-w /project \
		-e GOCACHE=/project/build/.gocache \
		-e GOMODCACHE=/project/build/.gomodcache \
		$(LINUX_BUILDER_IMAGE) \
		wails build -platform linux/amd64 -s -skipbindings -tags webkit2_41
	@echo "Linux binary ready: $(OUTPUT_DIR)/$(APP_NAME)"

# =============================================================================
#  Packaging & installers
# =============================================================================
.PHONY: package-macos package-windows package-linux
package-macos: ## Zip the macOS app bundle
	@if [ -d "$(OUTPUT_DIR)/$(BUNDLE_NAME)" ]; then \
		cd $(OUTPUT_DIR) && zip -r ../../$(APP_NAME)-macOS-universal.zip $(BUNDLE_NAME); \
	else \
		echo "No macOS bundle found. Run 'make build-macos' first."; exit 1; \
	fi

package-windows: ## Zip the Windows executable
	@if [ -f "$(OUTPUT_DIR)/$(APP_NAME)$(EXE_EXT)" ]; then \
		cd $(OUTPUT_DIR) && zip -r ../../$(APP_NAME)-windows-amd64.zip $(APP_NAME)$(EXE_EXT); \
	else \
		echo "No Windows binary found. Run 'make build-windows' first."; exit 1; \
	fi

package-linux: ## Zip the Linux executable
	@if [ -f "$(OUTPUT_DIR)/$(APP_NAME)" ]; then \
		cd $(OUTPUT_DIR) && zip -r ../../$(APP_NAME)-linux-amd64.zip $(APP_NAME); \
	else \
		echo "No Linux binary found. Run 'make build-linux' first."; exit 1; \
	fi

# --- Installers ---------------------------------------------------------------

NFPM := $(shell command -v nfpm 2>/dev/null)
ifeq ($(NFPM),)
  ifneq ($(wildcard build/.gopath/bin/nfpm),)
    NFPM := build/.gopath/bin/nfpm
  endif
endif
MAKENSIS := $(shell command -v makensis 2>/dev/null)

.PHONY: installer-macos dmg
installer-macos: dmg ## Alias for the macOS DMG target
dmg: ## Build .app, then package a drag-to-Applications DMG (macOS only)
	@if [ "$(OS_DETECTED)" != "macOS" ]; then \
		echo "ERROR: 'make dmg' must run on a macOS host (hdiutil is required)."; \
	echo "       Use CI (.github/workflows/release.yml) or a Mac instead."; exit 1; \
	fi
	@if [ ! -d "$(OUTPUT_DIR)/$(BUNDLE_NAME)" ]; then \
		echo "App bundle not found — building it first..."; \
		$(MAKE) --no-print-directory build-macos; \
	fi
	@rm -rf "$(DIST_DIR)/dmg-staging"
	mkdir -p "$(DIST_DIR)"
	# Stage: app + symlink to /Applications so users can drag-install.
	mkdir -p "$(DIST_DIR)/dmg-staging"
	cp -R "$(OUTPUT_DIR)/$(BUNDLE_NAME)" "$(DIST_DIR)/dmg-staging/"
	ln -sf /Applications "$(DIST_DIR)/dmg-staging/Applications"
	hdiutil create -volname "GetDevice $(VERSION)" \
		-srcfolder "$(DIST_DIR)/dmg-staging" \
		-ov -format UDZO \
		"$(DIST_DIR)/$(APP_NAME)-macOS-$(VERSION).dmg"
	rm -rf "$(DIST_DIR)/dmg-staging"
	@echo "DMG created: $(DIST_DIR)/$(APP_NAME)-macOS-$(VERSION).dmg"

.PHONY: sign-macos
sign-macos: ## Ad-hoc codesign the built .app (no Apple Developer ID needed)
	@if [ "$(OS_DETECTED)" != "macOS" ]; then \
		echo "ERROR: 'make sign-macos' must run on macOS."; exit 1; \
	fi
	@if [ ! -d "$(OUTPUT_DIR)/$(BUNDLE_NAME)" ]; then \
		echo "ERROR: no bundle at $(OUTPUT_DIR)/$(BUNDLE_NAME); run 'make build-macos' first."; exit 1; \
	fi
	codesign --force --deep --entitlements build/darwin/entitlements.plist \
		--sign - "$(OUTPUT_DIR)/$(BUNDLE_NAME)"
	@echo "Ad-hoc signed: $(OUTPUT_DIR)/$(BUNDLE_NAME)"

.PHONY: installer-windows nsis
installer-windows: nsis ## Alias for the Windows NSIS target
nsis: ## Build GetDevice.exe + NSIS setup installer (Windows host)
	@command -v wails >/dev/null 2>&1 || { echo "ERROR: Wails CLI not found. Run 'make deps' first."; exit 1; }
	@if [ -z "$(MAKENSIS)" ]; then \
		echo "WARNING: makensis not found on PATH; Wails will skip installer creation."; \
		echo "         Install NSIS 3+ and re-run: https://nsis.sourceforge.io"; \
	fi
	$(WAILS) build -platform windows/amd64 -s -skipbindings -nsis
	@echo "NSIS installer created in $(OUTPUT_DIR)/ (getdevice-api-amd64-installer.exe)"

.PHONY: nfpm-check
nfpm-check: ## Verify nfpm is installed
	@if [ -z "$(NFPM)" ]; then \
		echo "ERROR: nfpm is required to build Linux packages."; \
		echo "Install it with:"; \
		echo "  go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest"; \
		echo "(or let the Makefile use build/.gopath/bin/nfpm — run:"; \
		echo "  GOPATH=\"\$$PWD/build/.gopath\" go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest )"; \
		exit 1; \
	fi

.PHONY: deb rpm tarball linux-tarball
deb: ## Build the Linux .deb package (needs nfpm + a prior linux build)
	$(MAKE) --no-print-directory nfpm-check
	@if [ ! -f "$(OUTPUT_DIR)/$(APP_NAME)" ]; then \
		echo "Linux binary missing — building it first..."; \
		if [ "$(OS_DETECTED)" = "Linux" ]; then \
			$(MAKE) --no-print-directory build-linux; \
		else \
			$(MAKE) --no-print-directory build-linux-docker; \
		fi; \
	fi
	@mkdir -p "$(DIST_DIR)"
	VERSION="$(PKG_VERSION)" $(NFPM) package -f build/linux/nfpm.yaml -p deb \
		-t "$(DIST_DIR)/"

rpm: ## Build the Linux .rpm package (needs nfpm + a prior linux build)
	$(MAKE) --no-print-directory nfpm-check
	@if [ ! -f "$(OUTPUT_DIR)/$(APP_NAME)" ]; then \
		echo "Linux binary missing — building it first..."; \
		if [ "$(OS_DETECTED)" = "Linux" ]; then \
			$(MAKE) --no-print-directory build-linux; \
		else \
			$(MAKE) --no-print-directory build-linux-docker; \
		fi; \
	fi
	@mkdir -p "$(DIST_DIR)"
	VERSION="$(PKG_VERSION)" $(NFPM) package -f build/linux/nfpm.yaml -p rpm \
		-t "$(DIST_DIR)/"

linux-tarball: tarball ## Alias
tarball: ## Build a portable Linux tar.gz (needs a prior linux build)
	@if [ ! -f "$(OUTPUT_DIR)/$(APP_NAME)" ]; then \
		echo "Linux binary missing — building it first..."; \
		if [ "$(OS_DETECTED)" = "Linux" ]; then \
			$(MAKE) --no-print-directory build-linux; \
		else \
			$(MAKE) --no-print-directory build-linux-docker; \
		fi; \
	fi
	@mkdir -p "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64"
	cp "$(OUTPUT_DIR)/$(APP_NAME)" \
		"$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/$(APP_NAME)"
	cp README.md "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/" 2>/dev/null || true
	cp .env.example "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/.env.example"
	tar -czf "$(DIST_DIR)/$(APP_NAME)-linux-$(VERSION)-amd64.tar.gz" \
		-C "$(DIST_DIR)" "$(APP_NAME)-$(VERSION)-linux-amd64"
	rm -rf "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64"
	@echo "Tarball created: $(DIST_DIR)/$(APP_NAME)-linux-$(VERSION)-amd64.tar.gz"

.PHONY: installer-linux
installer-linux: deb rpm tarball ## Build all Linux packages (.deb, .rpm, .tar.gz)

.PHONY: installer
installer: ## Package an installer for the detected host OS
ifeq ($(OS_DETECTED),macOS)
	$(MAKE) --no-print-directory dmg
else ifeq ($(OS_DETECTED),Windows)
	$(MAKE) --no-print-directory nsis
else ifeq ($(OS_DETECTED),Linux)
	$(MAKE) --no-print-directory installer-linux
else
	@echo "ERROR: unsupported host OS '$(OS_DETECTED)'." && exit 1
endif

.PHONY: gen-icon
gen-icon: ## Regenerate the Linux PNG icon from build/appicon.png
	./build/linux/scripts/gen-icon.sh

# =============================================================================
#  Quality & maintenance
# =============================================================================
.PHONY: test
test: ## Run Go tests
	go test ./...

.PHONY: vet
vet: ## Run `go vet`
	go vet ./...

.PHONY: fmt
fmt: ## Format Go source files
	gofmt -w .

.PHONY: tidy
tidy: ## Tidy Go module dependencies
	go mod tidy

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(OUTPUT_DIR) $(DIST_DIR) *.zip
	@echo "Removed $(OUTPUT_DIR), $(DIST_DIR), and zip artifacts."
