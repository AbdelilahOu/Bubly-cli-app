# Makefile for Bubly CLI app

# Variables
BINARY_NAME=bubly
BUILD_DIR=build
SOURCE_FILES=main.go app/*.go utils/*.go types/*.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o ${BUILD_DIR}/${BINARY_NAME} main.go

# Run the application
.PHONY: run
run:
	make build-windows
	./build/bubly.exe

# Build for Windows
.PHONY: build-windows
build-windows:
	go build -o ${BUILD_DIR}/${BINARY_NAME}.exe main.go

# Build for Linux
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME} main.go

# Build for macOS
.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME} main.go

# Clean build directory
.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}

# Install dependencies
.PHONY: deps
deps:
	go mod tidy

# Run tests (if any)
.PHONY: test
test:
	go test ./...

# Run with verbose output
.PHONY: run-v
run-v:
	go run -v main.go

# Build and run
.PHONY: build-run
build-run: build
	./${BUILD_DIR}/${BINARY_NAME}

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  build-windows - Build for Windows"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-mac    - Build for macOS"
	@echo "  clean        - Clean build directory"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  run-v        - Run with verbose output"
	@echo "  build-run    - Build and run"
	@echo "  help         - Show this help"
