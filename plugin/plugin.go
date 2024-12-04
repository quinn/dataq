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
	Extract(context.Context, *pb.PluginRequest) (<-chan *pb.DataItem, error)
}

func GenerateHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Run is a harness that a plugin written in Go can use to receive and respond to requests
// plugins can be written in any language, as long as they implement the wire protocol
func Run(p Plugin) {
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
		// Create context for plugin execution
		ctx := context.Background()

		// Always configure the plugin first
		if err := p.Configure(req.Config); err != nil {
			stream.WriteResponse(os.Stdout, &pb.PluginResponse{
				PluginId: p.ID(),
				Error:    fmt.Sprintf("failed to configure plugin: %v", err),
			})
			continue
		}

		if req.Operation == "extract" {
			items, err := p.Extract(ctx, req)
			if err != nil {
				stream.WriteResponse(os.Stdout, &pb.PluginResponse{
					PluginId: p.ID(),
					Error:    fmt.Sprintf("failed to extract data: %v", err),
				})
				continue
			}

			for item := range items {
				// Infinite loop to allow debugger attachment
				for waitForDebugger {
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

			if req.GetClosed() {
				stream.WriteResponse(os.Stdout, &pb.PluginResponse{
					PluginId: p.ID(),
					Closed:   true,
				})

				return
			}

			// Signal end of this extract operation
			// stream.WriteResponse(os.Stdout, &pb.PluginResponse{
			// 	PluginId: p.ID(),
			// 	Done:    true,
			// })
			continue
		}

		stream.WriteResponse(os.Stdout, &pb.PluginResponse{
			PluginId: p.ID(),
			Error:    fmt.Sprintf("unknown operation: %s", req.Operation),
		})
	}
}
