package main

import (
	"context"
	"fmt"
	"io"
	"log"
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

	log.Printf("Received request: %s", string(input))

	// Parse request
	var req pb.PluginRequest
	if err := protojson.Unmarshal(input, &req); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing request: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Parsed request: %+v", req)

	// Process request
	var resp pb.PluginResponse
	resp.PluginId = plugin.ID()

	// Always configure the plugin first
	if err := plugin.Configure(req.Config); err != nil {
		resp.Error = err.Error()
	} else if req.Operation == "extract" {
		items, err := plugin.Extract(context.Background())
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

	// Write response
	output, err := protojson.Marshal(&resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
		os.Exit(1)
	}
}
