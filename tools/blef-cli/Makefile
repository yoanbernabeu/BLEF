.PHONY: build test install clean run-validate run-convert run-view help

# Binary name
BINARY_NAME=blef-cli

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOINSTALL)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	$(GOCMD) clean

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run validation on example file
run-validate:
	@echo "Running validate command on example..."
	./$(BINARY_NAME) validate ../../examples/minimal.blef.json

# Run view on example file
run-view:
	@echo "Running view command on example..."
	./$(BINARY_NAME) view ../../examples/complete.blef.json

# Build and run
run: build
	./$(BINARY_NAME) --help

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  install       - Install the binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  fmt           - Format code"
	@echo "  run-validate  - Run validate on example file"
	@echo "  run-view      - Run view on example file"
	@echo "  run           - Build and show help"
	@echo "  help          - Show this help message"

# Default target
.DEFAULT_GOAL := help

