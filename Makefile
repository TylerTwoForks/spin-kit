APP_NAME := spin-kit
CMD_PATH := ./cmd/spin-kit
DIST_DIR := dist

CGO_ENABLED ?= 0
LINUX_ARCH ?= amd64
MACOS_ARCH ?= arm64
WINDOWS_ARCH ?= amd64

.DEFAULT_GOAL := build

.PHONY: build build-linux build-macos build-windows clean help

build: build-linux build-macos build-windows
	@echo "Binaries are in $(DIST_DIR)/"

build-linux:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=$(LINUX_ARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(DIST_DIR)/$(APP_NAME)-linux-$(LINUX_ARCH) $(CMD_PATH)

build-macos:
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=$(MACOS_ARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(DIST_DIR)/$(APP_NAME)-macos-$(MACOS_ARCH) $(CMD_PATH)

build-windows:
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=$(WINDOWS_ARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(DIST_DIR)/$(APP_NAME)-windows-$(WINDOWS_ARCH).exe $(CMD_PATH)

clean:
	rm -rf $(DIST_DIR)

help:
	@echo "Targets:"
	@echo "  make                Build Linux/macOS/Windows binaries"
	@echo "  make build-linux    Build Linux binary"
	@echo "  make build-macos    Build macOS binary"
	@echo "  make build-windows  Build Windows binary"
	@echo "  make clean          Remove dist directory"
	@echo ""
	@echo "Optional arch overrides:"
	@echo "  make LINUX_ARCH=arm64 MACOS_ARCH=amd64 WINDOWS_ARCH=arm64"
