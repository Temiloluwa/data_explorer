.PHONY: all build clean test lint proto run-server run-cli help

# Variables
APP_NAME=data-explorer
CMD_SERVER_PATH=./cmd/server
CMD_CLI_PATH=./cmd/cli
OUTPUT_DIR=./bin
SERVER_BINARY=${OUTPUT_DIR}/${APP_NAME}-server
CLI_BINARY=${OUTPUT_DIR}/${APP_NAME}-cli
PROTO_API_DIR=./api/proto/dataexplorer/v1

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
# Ensure golangci-lint is installed: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GOLINT=golangci-lint
GOLINT_RUN=$(GOLINT) run ./... --timeout 5m
GOFMT=gofmt -w -s # Format code
GOVET=$(GOCMD) vet ./...

# Default target
all: build

# Build binaries
build: proto build-server build-cli

build-server:
	@echo "Building server binary..."
	@mkdir -p $(OUTPUT_DIR)
	$(GOBUILD) -o $(SERVER_BINARY) $(CMD_SERVER_PATH)/main.go

build-cli:
	@echo "Building CLI binary..."
	@mkdir -p $(OUTPUT_DIR)
	$(GOBUILD) -o $(CLI_BINARY) $(CMD_CLI_PATH)/main.go

# Generate protobuf/gRPC code
proto:
	@echo "Generating Protobuf/gRPC code..."
	@chmod +x ./scripts/gen_proto.sh
	@./scripts/gen_proto.sh
	@# Optional: Format generated code if needed
	# @$(GOFMT) ./api/proto

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@# Check if golangci-lint is installed
	@command -v $(GOLINT) >/dev/null 2>&1 || { echo >&2 "golangci-lint not found. Please install it: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	@$(GOLINT_RUN)

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOFMT) .

# Vet code
vet:
	@echo "Running go vet..."
	@$(GOVET)

# Clean build artifacts and generated code
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleaning generated proto code..."
	@rm -f $(PROTO_API_DIR)/*.pb.go
	@rm -f $(PROTO_API_DIR)/*.pb.gw.go
	@rm -f $(PROTO_API_DIR)/*.openapi.json
	# @$(GOCLEAN) # Optional: run go clean

# Run the server (requires building first)
run-server: build-server
	@echo "Running server..."
	$(SERVER_BINARY) # Add flags if needed, e.g., --config ./configs/config.yaml

# Run the CLI (requires building first)
run-cli: build-cli
	@echo "Running CLI (example: asking 'hello?')..."
	$(CLI_BINARY) ask "hello?" # Add flags if needed

# Help message
help:
	@echo "Available targets:"
	@echo "  all          - Build all binaries (default)"
	@echo "  build        - Build server and CLI binaries"
	@echo "  build-server - Build only the server binary"
	@echo "  build-cli    - Build only the CLI binary"
	@echo "  proto        - Generate Protobuf/gRPC code"
	@echo "  test         - Run unit tests"
	@echo "  lint         - Run the Go linter (golangci-lint)"
	@echo "  fmt          - Format Go code using gofmt"
	@echo "  vet          - Run go vet"
	@echo "  clean        - Remove build artifacts and generated code"
	@echo "  run-server   - Build and run the gRPC/REST server"
	@echo "  run-cli      - Build and run the CLI tool with an example command"
	@echo "  help         - Show this help message"
