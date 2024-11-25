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

# Clean generated files
clean:
    Remove-Item -Force -Recurse -ErrorAction SilentlyContinue proto/*.pb.go
    Remove-Item -Force -ErrorAction SilentlyContinue dataq.exe

# Setup everything (install tools, protoc, and generate code)
setup: install-tools install-protoc generate-proto

# Build the project
build:
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

# Build and run the example
run-example: build
    cd example
    dataq -config config.yaml
