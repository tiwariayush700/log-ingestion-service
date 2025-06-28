# Variables
APP_NAME := log-ingestion-service
DOCKER_IMAGE := yourusername/$(APP_NAME)
DOCKER_TAG := latest

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOPATH := $(shell go env GOPATH)
GOMODULE := $(shell go list -m)

# Build variables
BUILD_DIR := $(GOBASE)/build
BINARY_NAME := $(APP_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix

# Test variables
COVERAGE_DIR := $(GOBASE)/coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Dependency management
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(GOBIN)/$(BINARY_NAME) ./cmd

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_UNIX) ./cmd

# Run the application
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	go run ./cmd

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(GOBIN)
	rm -rf $(COVERAGE_DIR)
	rm -rf $(BUILD_DIR)

# Test
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Test with short flag (skips integration tests)
.PHONY: test-short
test-short:
	@echo "Running short tests..."
	go test -v -short ./...

# Test with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Generate test coverage
.PHONY: test-coverage
test-coverage:
	@echo "Generating test coverage..."
	mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

# Lint
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up
.PHONY: docker-up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build -d

# Docker compose down
.PHONY: docker-down
docker-down:
	@echo "Stopping services with Docker Compose..."
	docker-compose down

# Docker push
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

# Generate Go module documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	mkdir -p docs
	godoc -http=:6060 &
	@echo "Documentation server running at http://localhost:6060/pkg/$(GOMODULE)/"

# Setup development environment
.PHONY: setup
setup: deps
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest

# All-in-one development setup
.PHONY: dev-setup
dev-setup: clean setup fmt vet lint test-coverage

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  deps            - Install dependencies"
	@echo "  build           - Build the application"
	@echo "  build-linux     - Build for Linux"
	@echo "  run             - Run the application"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  test-short      - Run tests with short flag (skips integration tests)"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-coverage   - Generate test coverage report"
	@echo "  lint            - Lint code"
	@echo "  fmt             - Format code"
	@echo "  vet             - Vet code"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-up       - Start services with Docker Compose"
	@echo "  docker-down     - Stop services with Docker Compose"
	@echo "  docker-push     - Push Docker image"
	@echo "  docs            - Generate Go module documentation"
	@echo "  setup           - Setup development environment"
	@echo "  dev-setup       - All-in-one development setup"
	@echo "  help            - Show this help message"

# Default target
.DEFAULT_GOAL := help