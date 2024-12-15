package pluginutil

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/hash"
	"go.quinn.io/dataq/schema"
	"go.quinn.io/dataq/stream"
)

// ExecutePlugin executes a plugin binary with the given request and returns a channel of responses
func Execute(ctx context.Context, cfg *config.Plugin, req <-chan *schema.PluginRequest) (<-chan *schema.PluginResponse, error) {
	binpath := filepath.Join(config.StateDir(), "bin", cfg.BinaryPath)
	if _, err := os.Stat(binpath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", binpath)
	}

	pluginDir := filepath.Join(config.StateDir(), cfg.ID)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, binpath)
	cmd.Dir = pluginDir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	// Start the plugin
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin: %w", err)
	}

	// Create response channel
	responses := make(chan *schema.PluginResponse)

	go func() {
		errBytes, _ := io.ReadAll(stderr)
		responses <- &schema.PluginResponse{
			PluginId: cfg.ID,
			Error:    "plugin stderr: " + string(errBytes),
		}

		responses <- &schema.PluginResponse{
			Closed: true,
		}
	}()

	go func() {
		for req := range req {
			stream.WriteRequest(stdin, req)
		}
	}()

	stream, errc := stream.StreamResponses(stdout)

	// Handle plugin execution in goroutine
	go func() {
		// Stream responses until error or EOF
		for resp := range stream {
			if resp == nil {
				// EOF, probably. Seems ok, check errc
				continue
			}

			if resp.Item != nil {
				if resp.Item.Meta.Hash != hash.Generate(resp.Item.RawData) {
					responses <- &schema.PluginResponse{
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
	}()

	go func() {
		// Check for stream error
		if err := <-errc; err != nil {
			responses <- &schema.PluginResponse{
				PluginId: cfg.ID,
				Error:    fmt.Sprintf("stream error: %v", err),
			}

			responses <- &schema.PluginResponse{
				Closed: true,
			}

			return
		}
	}()

	go func() {
		// Wait for plugin to finish
		if err := cmd.Wait(); err != nil {
			responses <- &schema.PluginResponse{
				PluginId: cfg.ID,
				Error:    fmt.Sprintf("plugin execution failed: %v\n", err),
			}

			responses <- &schema.PluginResponse{
				Closed: true,
			}
		}
	}()

	return responses, nil
}
