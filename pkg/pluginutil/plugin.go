package pluginutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/stream"
	"google.golang.org/protobuf/proto"
)

// Plugin is the interface that all plugins must implement
type Plugin interface {
	ID() string
	Configure(map[string]string) error
	Extract(context.Context, *pb.PluginRequest) (<-chan *pb.DataItem, error)
}

// ExecutePlugin executes a plugin binary with the given request and returns a channel of responses
func ExecutePlugin(ctx context.Context, plugin *plugin.PluginConfig, data *pb.DataItem) (<-chan *pb.PluginResponse, error) {
	if _, err := os.Stat(plugin.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", plugin.BinaryPath)
	}

	// Create plugin request
	req := &pb.PluginRequest{
		PluginId:  plugin.ID,
		Operation: "extract",
		Config:    plugin.Config,
		Item:      data,
	}

	// Serialize request using protobuf
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, plugin.BinaryPath)
	cmd.Stdin = bytes.NewReader(reqData)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Start the plugin
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin: %w", err)
	}

	// Create response channel
	responses := make(chan *pb.PluginResponse)

	// Handle plugin execution in goroutine
	go func() {
		defer close(responses)

		// Stream responses until error or EOF
		stream, errc := stream.StreamResponses(stdout)
		for resp := range stream {
			responses <- resp
		}

		// Check for stream error
		if err := <-errc; err != nil {
			responses <- &pb.PluginResponse{
				PluginId: plugin.ID,
				Error:    fmt.Sprintf("stream error: %v", err),
			}
			return
		}

		// Wait for plugin to finish
		if err := cmd.Wait(); err != nil {
			responses <- &pb.PluginResponse{
				PluginId: plugin.ID,
				Error:    fmt.Sprintf("plugin execution failed: %v\nstderr: %s", err, stderr.String()),
			}
		}
	}()

	return responses, nil
}

// HandlePlugin handles plugin execution and writes responses to stdout
func HandlePlugin(p Plugin) error {
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
