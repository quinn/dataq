# DataQ

DataQ is a general-purpose tool for archiving and organizing data from various sources using a plugin-based architecture.

## Setup

1. Install the Protocol Buffer compiler (protoc):
   - Windows: Download from https://github.com/protocolbuffers/protobuf/releases
   - Add protoc to your PATH

2. Install Go dependencies:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go mod download
```

3. Generate Protocol Buffer code:
```bash
protoc --go_out=. --go_opt=paths=source_relative proto/dataq.proto
```

## Project Structure

- `proto/`: Protocol Buffer definitions
- `plugin/`: Plugin interface and registry
- `plugin/filescan/`: File system scanner plugin implementation

## Usage

1. Build the project:
```bash
go build
```

2. Run the example:
```bash
./dataq
```

The example will scan the current directory and display information about all files found.

## Creating New Plugins

To create a new plugin:

1. Implement the Plugin interface defined in `plugin/plugin.go`
2. Register your plugin with the registry in your application

Example plugin implementation can be found in `plugin/filescan/filescan.go`.
