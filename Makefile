.PHONY: all build clean test lint deps

# Default target
all: deps test lint build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go install mvdan.cc/garble@latest
	@if ! command -v upx &> /dev/null; then echo "Warning: UPX not installed. Binary compression will be skipped."; fi

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run

# Build for all platforms
build:
	@chmod +x scripts/build.sh
	@./scripts/build.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/

# Build for current platform (development)
build-dev:
	@echo "Building for development..."
	@go build -o build/apex . 