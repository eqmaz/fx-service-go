# Makefile

# Suppress command echoing globally
.SILENT:

# Variables
BINARY_NAME=server
BUILD_DIR=.
MAIN_SRC=cmd/server/main.go

# Help target
.PHONY: help
help:
	echo "Usage:"
	echo "  make build      - Compile the Go project and output the binary"
	echo "  make min-build  - Compile the Go project with -ldflags=\"-s -w\" and output the binary"
	echo "  make run        - Build and then run the project"
	echo "  make clean      - Remove the compiled binary"
	echo "  make start      - Start (run) the already compiled binary"
	echo "  make test       - Run tests"
	echo "  make test-cover - Run tests with coverage"

# Default target
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Build target
.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_SRC)

# Minimal build target
.PHONY: min-build
min-build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_SRC)

# Clean target
.PHONY: clean
clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

# Start target
.PHONY: start
start:
	./$(BUILD_DIR)/$(BINARY_NAME)

# Test target
.PHONY: test
test:
	go test -v ./...

# Test coverage target
.PHONY: test-cover
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
