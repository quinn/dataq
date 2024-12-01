package pluginutil

import (
	"context"
	"fmt"
	"io"
	"os"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/proto"
)

// Plugin is the interface that all plugins must implement
type Plugin interface {
	ID() string
	Configure(map[string]string) error
	Extract(context.Context, *pb.PluginRequest) (<-chan *pb.DataItem, error)
}

// HandlePlugin reads the plugin request from stdin, processes it using the given plugin,
// and writes the response to stdout
func HandlePlugin(p Plugin) {
	// Read request from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Parse request
	var req pb.PluginRequest
	if err := proto.Unmarshal(input, &req); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing request: %v\n", err)
		os.Exit(1)
	}

	// Process request
	var resp pb.PluginResponse
	resp.PluginId = p.ID()

	// Always configure the plugin first
	if err := p.Configure(req.Config); err != nil {
		resp.Error = err.Error()
	} else if req.Operation == "extract" {
		items, err := p.Extract(context.Background(), &req)
		if err != nil {
			resp.Error = err.Error()
		} else {
			for item := range items {
				resp.Items = append(resp.Items, item)
			}
		}
	} else if req.Operation != "configure" {
		resp.Error = fmt.Sprintf("unknown operation: %s", req.Operation)
	}

	// Write response using protobuf
	output, err := proto.Marshal(&resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
		os.Exit(1)
	}
}
