# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/cli/main.go

# Run the application
run:
	@go run cmd/cli/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	# @go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run test clean itest
