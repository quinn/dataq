# protoc-gen-dataq-index

A protoc plugin that generates Go code for indexing protobuf messages in DataQ.

## Installation

```bash
go install go.quinn.io/dataq/cmd/protoc-gen-dataq-index@latest
```

## Usage

1. Add the following option to your .proto files to enable index generation:
```protobuf
option go_plugin_indexable = true;
```

2. Run protoc with the dataq-index plugin:
```bash
protoc --go_out=. \
       --go-grpc_out=. \
       --dataq-index_out=. \
       *.proto
```

This will generate `*_index.pb.go` files containing the following methods for each message:

- `SchemaKind() string`: Returns the message name as the schema kind
- `Metadata() map[string]interface{}`: Returns a map of all non-zero fields in the message

## Example

For a protobuf message:
```protobuf
message ExtractRequest {
    string plugin_id = 1;
    string hash = 2;
    string kind = 3;
    map<string, string> metadata = 4;
}
```

The plugin will generate:
```go
func (m *ExtractRequest) SchemaKind() string {
    return "ExtractRequest"
}

func (m *ExtractRequest) Metadata() map[string]interface{} {
    metadata := make(map[string]interface{})
    if m.PluginId != "" {
        metadata["plugin_id"] = m.PluginId
    }
    if m.Hash != "" {
        metadata["hash"] = m.Hash
    }
    if m.Kind != "" {
        metadata["kind"] = m.Kind
    }
    if m.Metadata != nil {
        metadata["metadata"] = m.Metadata
    }
    return metadata
}
```
