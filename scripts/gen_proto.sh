#!/bin/bash

# Exit on error
set -e

# --- Configuration ---
# Automatically determine the module path from go.mod
MODULE_PATH=$(go list -m)
if [ -z "$MODULE_PATH" ]; then
    echo "Error: Could not determine Go module path. Run 'go mod init <your_module_path>' first."
    exit 1
fi
echo "Using Go module path: $MODULE_PATH"

PROTO_DIR=./api/proto
PROTO_API_DIR=${PROTO_DIR}/dataexplorer/v1
GEN_DIR=./api/proto # Output generated Go code alongside proto files
# --- End Configuration ---


# Ensure protoc is installed
if ! command -v protoc &> /dev/null
then
    echo "protoc could not be found. Please install Protocol Buffers compiler."
    echo "See: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Ensure Go plugins are installed
echo "Checking for Go gRPC/Protobuf/Gateway plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest # Optional: for OpenAPI spec

# Ensure GOPATH/bin is in the PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Create output directory if it doesn't exist
mkdir -p ${GEN_DIR}

# Define the Go package option based on the module path
GO_PACKAGE_OPT="${MODULE_PATH}/api/proto/dataexplorer/v1;dataexplorerv1"

# Find the proto file
PROTO_FILE=$(find ${PROTO_API_DIR} -maxdepth 1 -name '*.proto' | head -n 1)
if [ -z "$PROTO_FILE" ]; then
    echo "Error: No .proto file found in ${PROTO_API_DIR}"
    exit 1
fi
echo "Found proto file: $PROTO_FILE"

# Generate Go gRPC and Protobuf code
echo "Generating Go gRPC code..."
protoc --proto_path=${PROTO_DIR} \
       --go_out=${GEN_DIR} --go_opt=paths=source_relative \
       --go-grpc_out=${GEN_DIR} --go-grpc_opt=paths=source_relative \
       ${PROTO_FILE}

# Generate gRPC Gateway code (REST proxy)
echo "Generating gRPC Gateway code..."
protoc --proto_path=${PROTO_DIR} \
       --grpc-gateway_out=${GEN_DIR} \
       --grpc-gateway_opt paths=source_relative \
       --grpc-gateway_opt generate_unbound_methods=true \
       ${PROTO_FILE}

# Optional: Generate OpenAPIv2 specification
echo "Generating OpenAPIv2 specification..."
protoc --proto_path=${PROTO_DIR} \
      --openapiv2_out=${GEN_DIR} \
      --openapiv2_opt logtostderr=true \
      ${PROTO_FILE}


echo "Code generation complete."

# Note: This script assumes googleapis protos needed for annotations are available.
# If you encounter errors like "google/api/annotations.proto not found", you might need to:
# 1. Vendor them using a tool like Buf (https://buf.build)
# 2. Clone the googleapis repository (e.g., into a 'third_party' directory) and add
#    `--proto_path=./third_party/googleapis` to the protoc commands above.
