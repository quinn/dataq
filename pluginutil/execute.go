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
)

// ExecutePlugin executes a plugin binary with the given request and returns a channel of responses
func Execute(ctx context.Context, cfg *config.Plugin, req chan *pb.PluginRequest) (<-chan *pb.PluginResponse, error) {
	if _, err := os.Stat(cfg.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", cfg.BinaryPath)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, cfg.BinaryPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Start the plugin
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin: %w", err)
	}

	// Create response channel
	responses := make(chan *pb.PluginResponse)

	for req := range req {
		stream.WriteRequest(stdin, req)

		// Handle plugin execution in goroutine
		go func() {
			// Stream responses until error or EOF
			stream, errc := stream.StreamResponses(stdout)
			for resp := range stream {
				if resp == nil {
					// EOF, probably. Seems ok, check errc
					continue
				}

				if resp.Item != nil {
					if resp.Item.Meta.Hash != plugin.GenerateHash(resp.Item.RawData) {
						responses <- &pb.PluginResponse{
							PluginId: cfg.ID,
							Error:    "hash mismatch",
						}
						return
					}
				}

				// Closed message body is discarded
				if resp.GetClosed() {
					close(responses)
					return
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
	}

	return responses, nil
}
