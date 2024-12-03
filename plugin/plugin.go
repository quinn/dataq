package plugin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"go.quinn.io/dataq/stream"
	"google.golang.org/protobuf/proto"

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

		return
	}

	stream.WriteResponse(os.Stdout, &pb.PluginResponse{
		PluginId: p.ID(),
		Error:    fmt.Sprintf("unknown operation: %s", req.Operation),
	})
}
