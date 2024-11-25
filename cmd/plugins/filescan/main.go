package main

import (
	"context"
	"fmt"
	"io"
	"os"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	plugin := New()

	// Read request from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Parse request
	var req pb.PluginRequest
	if err := protojson.Unmarshal(input, &req); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing request: %v\n", err)
		os.Exit(1)
	}

	// Process request
	var resp pb.PluginResponse
	resp.PluginId = plugin.ID()

	switch req.Operation {
	case "configure":
		if err := plugin.Configure(req.Config); err != nil {
			resp.Error = err.Error()
		}
	case "extract":
		items, err := plugin.Extract(context.Background())
		if err != nil {
			resp.Error = err.Error()
		} else {
			// Collect items from channel
			for item := range items {
				resp.Items = append(resp.Items, item)
			}
		}
	default:
		resp.Error = fmt.Sprintf("unknown operation: %s", req.Operation)
	}

	// Write response to stdout
	output, err := protojson.Marshal(&resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(output)
}
