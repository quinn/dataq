package pluginutil

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/stream"
	"google.golang.org/protobuf/proto"
)

// debuggerAttached is used as a breakpoint for debugger attachment
var debuggerAttached bool

// Plugin is the interface that all plugins must implement
type Plugin interface {
	ID() string
	Configure(map[string]string) error
	Extract(context.Context, *pb.PluginRequest) (<-chan *pb.DataItem, error)
}

// ExecutePlugin executes a plugin binary with the given request and returns a channel of responses
func ExecutePlugin(ctx context.Context, pc *plugin.PluginConfig, data *pb.DataItem) (<-chan *pb.PluginResponse, error) {
	if _, err := os.Stat(pc.BinaryPath); err != nil {
		return nil, fmt.Errorf("plugin binary not found: %s", pc.BinaryPath)
	}

	// Create plugin request
	req := &pb.PluginRequest{
		PluginId:  pc.ID,
		Operation: "extract",
		Config:    pc.Config,
		Item:      data,
	}

	// Serialize request using protobuf
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, pc.BinaryPath)
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
				if resp.Item.Meta.Hash != GenerateHash(resp.Item.RawData) {
					responses <- &pb.PluginResponse{
						PluginId: pc.ID,
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
				PluginId: pc.ID,
				Error:    fmt.Sprintf("stream error: %v", err),
			}
			return
		}

		// Wait for plugin to finish
		if err := cmd.Wait(); err != nil {
			responses <- &pb.PluginResponse{
				PluginId: pc.ID,
				Error:    fmt.Sprintf("plugin execution failed: %v\nstderr: %s", err, stderr.String()),
			}
		}
	}()

	return responses, nil
}

func GenerateHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// HandlePlugin handles plugin execution and writes responses to stdout
func HandlePlugin(p Plugin) {
	// Read request from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		stream.WriteResponse(os.Stdout, &pb.PluginResponse{
			PluginId: p.ID(),
			Error:    fmt.Sprintf("failed to read request: %v", err),
		})
		return
	}

	// Parse request using protobuf
	var req pb.PluginRequest
	if err := proto.Unmarshal(input, &req); err != nil {
		stream.WriteResponse(os.Stdout, &pb.PluginResponse{
			PluginId: p.ID(),
			Error:    fmt.Sprintf("failed to unmarshal request: %v", err),
		})
		return
	}

	// Create context for plugin execution
	ctx := context.Background()

	// Always configure the plugin first
	if err := p.Configure(req.Config); err != nil {
		stream.WriteResponse(os.Stdout, &pb.PluginResponse{
			PluginId: p.ID(),
			Error:    fmt.Sprintf("failed to configure plugin: %v", err),
		})
		return
	} else if req.Operation == "extract" {
		items, err := p.Extract(ctx, &req)
		if err != nil {
			stream.WriteResponse(os.Stdout, &pb.PluginResponse{
				PluginId: p.ID(),
				Error:    fmt.Sprintf("failed to extract data: %v", err),
			})
			return
		}

		for item := range items {
			// Infinite loop to allow debugger attachment
			debuggerAttached = false // Will be set to true via debugger
			for !debuggerAttached {
				time.Sleep(time.Millisecond) // Sleep briefly to prevent CPU spinning
			}
			item.Meta.Hash = GenerateHash(item.RawData)
			if req.Item != nil {
				item.Meta.ParentHash = req.Item.Meta.Hash
			}

			stream.WriteResponse(os.Stdout, &pb.PluginResponse{
				PluginId: item.Meta.PluginId,
				Item:     item,
			})
		}

		return
	}

	stream.WriteResponse(os.Stdout, &pb.PluginResponse{
		PluginId: p.ID(),
		Error:    fmt.Sprintf("unknown operation: %s", req.Operation),
	})
}
