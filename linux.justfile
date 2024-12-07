# Default recipe to run when just is called without arguments
default:
    @just --list

# Install required Go tools
install-tools:
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Download and install protoc for Linux
install-protoc:
    echo "Downloading protoc for Linux..."
    mkdir -p /tmp/protoc
    wget -P /tmp/protoc https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-linux-x86_64.zip
    unzip /tmp/protoc/protoc-25.1-linux-x86_64.zip -d /tmp/protoc
    sudo mv /tmp/protoc/bin/protoc $GOPATH/bin/
    rm -rf /tmp/protoc
    echo "protoc installed successfully"

# Generate Go code from proto files
generate-proto:
    protoc --go_out=. --go_opt=paths=source_relative proto/dataq.proto

# Build all plugins
build-plugins:
    echo "Building plugins..."
    mkdir -p cmd/plugins/filescan/bin
    mkdir -p cmd/plugins/gmail/bin
    go build -o $HOME/.config/dataq/state/bin/filescan cmd/plugins/filescan/main.go cmd/plugins/filescan/filescan.go
    go build -o $HOME/.config/dataq/state/bin/gmail cmd/plugins/gmail/main.go cmd/plugins/gmail/gmail.go
    echo "Plugins built successfully"

# Clean plugins
clean-plugins:
    echo "Cleaning plugins..."
    rm -rf cmd/plugins/*/bin

# Clean generated files and plugins
clean: clean-plugins
    rm -f proto/*.pb.go
    rm -f dataq

# Setup everything (install tools, protoc, and generate code)
setup: install-tools install-protoc generate-proto

# Build the project and plugins
build: build-plugins
    go build -o dataq cmd/main.go

# Run tests
test:
    go test ./...

# Format code
fmt:
    go fmt ./...

# Run linter
lint:
    go vet ./...

# Run example with plugins
run mode: build
    echo "Running dataq..."
    ./dataq --mode {{ mode }}
