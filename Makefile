# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o zoop main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	# @go test ./internal/database -v

# Cross-platform builds
build-linux-amd64:
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -o zoop-linux-amd64 main.go

build-darwin-amd64:
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build -o zoop-darwin-amd64 main.go

build-darwin-arm64:
	@echo "Building for macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 go build -o zoop-darwin-arm64 main.go

build-windows-amd64:
	@echo "Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build -o zoop-windows-amd64.exe main.go

build-all: build-linux-amd64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main zoop zoop-*

.PHONY: all build run test clean itest build-linux-amd64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-all
