package plugin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"go.quinn.io/dataq/stream"

	pb "go.quinn.io/dataq/proto"
)

// waitForDebugger is used as a breakpoint for debugger attachment
var waitForDebugger bool = false

// Plugin is the interface that all plugins must implement
type Plugin interface {
	ID() string
	Configure(map[string]string) error
	Extract(context.Context, *pb.PluginRequest, *PluginAPI) (<-chan *pb.DataItem, error)
}

// Run is a harness that a plugin written in Go can use to receive and respond to requests
// plugins can be written in any language, as long as they implement the wire protocol
func Run(p Plugin) {
	// Infinite loop to allow debugger attachment
	for waitForDebugger {
		time.Sleep(time.Millisecond) // Sleep briefly to prevent CPU spinning
	}

	requests, errc := stream.StreamRequests(os.Stdin)

	// Handle stream errors in a separate goroutine
	go func() {
		if err := <-errc; err != nil {
			stream.WriteResponse(os.Stdout, &pb.PluginResponse{
				PluginId: p.ID(),
				Error:    fmt.Sprintf("stream error: %v", err),
			})
		}
	}()

	for req := range requests {
		api := &PluginAPI{plugin: p, request: req}
		// Create context for plugin execution
		ctx := context.Background()

		// Always configure the plugin first
		if err := p.Configure(req.Config); err != nil {
			api.WriteError(fmt.Errorf("failed to configure plugin: %v", err))
			continue
		}

		// if req.Operation == "extract" {
		items, err := p.Extract(ctx, req, api)
		if err != nil {
			api.WriteError(fmt.Errorf("failed to extract data: %v", err))
			continue
		}

		for item := range items {
			api.WriteItem(item)
		}

		if req.GetClosed() {
			api.WriteClosed()
			return
		}

		// Signal end of this extract operation
		api.WriteDone()
	}
}
