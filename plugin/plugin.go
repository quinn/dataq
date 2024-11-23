package plugin

import (
	"context"

	pb "go.quinn.io/dataq/proto"
)

// Plugin defines the interface that all data extraction plugins must implement
type Plugin interface {
	// ID returns the unique identifier of the plugin
	ID() string

	// Name returns the human-readable name of the plugin
	Name() string

	// Description returns a description of what the plugin does
	Description() string

	// Configure sets up the plugin with the provided configuration
	Configure(config map[string]string) error

	// Extract performs the data extraction and returns a channel of DataItems
	Extract(ctx context.Context) (<-chan *pb.DataItem, error)
}
