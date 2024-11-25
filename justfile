set windows-powershell

# Default recipe to run when just is called without arguments
default:
    @just --list

# Install required Go tools
install-tools:
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Download and install protoc for Windows
install-protoc:
    Write-Host "Downloading protoc for Windows..." -ForegroundColor Green
    New-Item -ItemType Directory -Force -Path "$env:TEMP\protoc"
    Invoke-WebRequest -Uri 'https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-win64.zip' -OutFile "$env:TEMP\protoc\protoc.zip"
    Expand-Archive -Force "$env:TEMP\protoc\protoc.zip" -DestinationPath "$env:TEMP\protoc"
    Copy-Item "$env:TEMP\protoc\bin\protoc.exe" -Destination "$env:GOPATH\bin" -Force
    Remove-Item -Recurse -Force "$env:TEMP\protoc"
    Write-Host "protoc installed successfully" -ForegroundColor Green

# Generate Go code from proto files
generate-proto:
    protoc --go_out=. --go_opt=paths=source_relative proto/dataq.proto

# Build all plugins
build-plugins:
    Write-Host "Building plugins..." -ForegroundColor Green
    New-Item -ItemType Directory -Force -Path "cmd/plugins/filescan/bin"
    New-Item -ItemType Directory -Force -Path "cmd/plugins/gmail/bin"
    go build -o cmd/plugins/filescan/bin/filescan.exe cmd/plugins/filescan/main.go cmd/plugins/filescan/filescan.go
    go build -o cmd/plugins/gmail/bin/gmail.exe cmd/plugins/gmail/main.go cmd/plugins/gmail/gmail.go
    Write-Host "Plugins built successfully" -ForegroundColor Green

# Clean plugins
clean-plugins:
    Write-Host "Cleaning plugins..." -ForegroundColor Green
    Remove-Item -Force -Recurse -ErrorAction SilentlyContinue cmd/plugins/*/bin

# Clean generated files and plugins
clean: clean-plugins
    Remove-Item -Force -Recurse -ErrorAction SilentlyContinue proto/*.pb.go
    Remove-Item -Force -ErrorAction SilentlyContinue dataq.exe

# Setup everything (install tools, protoc, and generate code)
setup: install-tools install-protoc generate-proto

# Build the project and plugins
build: build-plugins
    go build -o dataq.exe

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
run-example: build
    Write-Host "Running example..." -ForegroundColor Green
    New-Item -ItemType Directory -Force -Path "example/plugins/filescan"
    Copy-Item -Force cmd/plugins/filescan/bin/filescan.exe example/plugins/filescan/
    cd example
    ../dataq -config config.yaml
