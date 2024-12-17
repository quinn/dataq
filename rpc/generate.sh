#!/bin/bash
set -e

# Install protoc Go plugins if not already installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Ensure the output directories exist
mkdir -p gen/go

# Generate Go code
protoc \
    --go_out=gen/go \
    --go_opt=paths=source_relative \
    --go-grpc_out=gen/go \
    --go-grpc_opt=paths=source_relative \
    --proto_path=. \
    schema.proto dataq.proto

# Update go.mod for the generated code
cd gen/go
go mod init github.com/quinn/dataq/rpc/gen/go
go mod tidy

echo "Protocol buffer code generation complete!"
