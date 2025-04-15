# Data Explorer

This project implements a Question Answering system called Data Explorer, primarily in Go.
It aims to answer user questions based on information retrieved from various sources like
web searches, local documents, databases, and external APIs.

The system exposes its core logic via gRPC, REST (using grpc-gateway), and a Command-Line Interface (CLI).

## Project Structure

(See previous description or generate using `tree`)

## Prerequisites

- Go (version 1.20 or later recommended)
- Protocol Buffers Compiler (`protoc`)
- Go gRPC/Protobuf plugins (`protoc-gen-go`, `protoc-gen-go-grpc`)
- gRPC Gateway plugins (`protoc-gen-grpc-gateway`, `protoc-gen-openapiv2`)

You can install the Go plugins using:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install [github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest)
go install [github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest)
```
Ensure `$(go env GOPATH)/bin` is in your system's `PATH`.

## Getting Started

1.  **Clone the repository (or use this bootstrap script):**
    ```bash
    # If you used the bootstrap script, you are already inside the project dir
    # git clone <your-repo-url>
    # cd data-explorer
    ```

2.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Generate gRPC/Protobuf Code:**
    ```bash
    make proto
    ```
    *Note: This requires `protoc` and the necessary plugins to be installed and in your PATH.*

4.  **Configure:**
    - Copy `configs/config.yaml` or create your own.
    - Set necessary environment variables (e.g., `DATAEXPLORER_DATASOURCES_WEB_SEARCH_API_KEY`). See `internal/config/config.go` and `configs/config.yaml` for details.

5.  **Build:**
    ```bash
    make build
    ```
    This creates `./bin/data-explorer-server` and `./bin/data-explorer-cli`.

6.  **Run:**
    * **Server (gRPC & REST):**
        ```bash
        make run-server
        # or
        ./bin/data-explorer-server --config ./configs/config.yaml
        ```
        The server will listen on ports defined in the config (default gRPC :50051, REST :8080/:8081).

    * **CLI:**
        ```bash
        make run-cli
        # or
        ./bin/data-explorer-cli ask "What is the capital of Nigeria?" --config ./configs/config.yaml
        ```
        Use `./bin/data-explorer-cli --help` for more options.

## Development

- **Linting:** `make lint` (Requires `golangci-lint`)
- **Formatting:** `make fmt`
- **Testing:** `make test`
- **Cleaning:** `make clean`

## TODO

- Implement actual data source fetching logic.
- Implement core QA engine logic (NLP/LLM integration).
- Add robust error handling and observability (metrics, tracing).
- Secure endpoints if necessary.
- Write comprehensive tests.
- Add deployment configurations (Dockerfile, etc.).
