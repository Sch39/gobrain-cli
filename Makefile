# =========================
# GoBrain Makefile
# =========================

APP_NAME := gob
CMD_PATH := ./cmd/gob
BIN_DIR  := bin

# OS-aware binary name (adds .exe on Windows)
EXE :=
ifeq ($(OS),Windows_NT)
	EXE := .exe
endif
BIN_PATH := $(BIN_DIR)/$(APP_NAME)$(EXE)

SANDBOX_DIR := sandbox

GO := go
GOFLAGS := -trimpath

# -------------------------
# Default
# -------------------------
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make sandbox   - Build gob & prepare sandbox for manual testing"
	@echo "  make test      - Run integration tests (tests/)"
	@echo "  make build     - Build local binary"
	@echo "  make release   - Build release binary"
	@echo "  make clean     - Remove build & sandbox artifacts"

# -------------------------
# Build
# -------------------------
.PHONY: build
build:
	@echo ">> Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_PATH) $(CMD_PATH)

# -------------------------
# Sandbox (Manual Test)
# -------------------------
.PHONY: sandbox
sandbox: build
	@echo ">> Preparing sandbox environment..."
	@rm -rf $(SANDBOX_DIR)
	@mkdir -p $(SANDBOX_DIR)
	@echo ""
	@echo "Sandbox ready!"
	@echo "-----------------------------------"
	@echo "Binary : $(CURDIR)/$(BIN_PATH)"
	@echo "Workdir: $(CURDIR)/$(SANDBOX_DIR)"
	@echo ""
	@echo "Example:"
	@echo "  cd $(SANDBOX_DIR) && ../$(BIN_PATH) init"
	@echo ""

# -------------------------
# Test (Integration)
# -------------------------
.PHONY: test
test: build
	@echo ">> Running integration tests..."
	$(GO) test ./tests -v

# -------------------------
# Release Build
# -------------------------
.PHONY: release
release:
	@echo ">> Building release binary..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 $(GO) build \
		$(GOFLAGS) \
		-ldflags "-s -w" \
		-o $(BIN_PATH) \
		$(CMD_PATH)
	@echo "Release binary at $(BIN_PATH)"

# -------------------------
# Clean
# -------------------------
.PHONY: clean
clean:
	@echo ">> Cleaning artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(SANDBOX_DIR)
