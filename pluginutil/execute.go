package pluginutil

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/stream"
	"google.golang.org/protobuf/proto"
)

// ExecutePlugin executes a plugin binary with the given request and returns a channel of responses
func Execute(ctx context.Context, cfg *config.Plugin, data *pb.DataItem) (<-chan *pb.PluginResponse, error) {
	if _, err := os.Stat(cfg.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", cfg.BinaryPath)
	}

	// Create plugin request
	req := &pb.PluginRequest{
		PluginId:  cfg.ID,
		Operation: "extract",
		Config:    cfg.Config,
		Item:      data,
	}

	// Serialize request using protobuf
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, cfg.BinaryPath)
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
			if resp.Item != nil {
				if resp.Item.Meta.Hash != plugin.GenerateHash(resp.Item.RawData) {
					responses <- &pb.PluginResponse{
						PluginId: cfg.ID,
						Error:    "hash mismatch",
					}
					return
				}
			}

			responses <- resp
		}

		// Check for stream error
		if err := <-errc; err != nil {
			responses <- &pb.PluginResponse{
				PluginId: cfg.ID,
				Error:    fmt.Sprintf("stream error: %v", err),
			}
			return
		}

		// Wait for plugin to finish
		if err := cmd.Wait(); err != nil {
			responses <- &pb.PluginResponse{
				PluginId: cfg.ID,
				Error:    fmt.Sprintf("plugin execution failed: %v\nstderr: %s", err, stderr.String()),
			}
		}
	}()

	return responses, nil
}
