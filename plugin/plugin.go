package plugin

import (
	"context"

	pb "go.quinn.io/dataq/proto"
)

// PluginConfig contains configuration for a plugin
type PluginConfig struct {
	ID         string            `yaml:"id"`
	Enabled    bool              `yaml:"enabled"`
	BinaryPath string            `yaml:"binary_path"`
	Config     map[string]string `yaml:"config"`
}

// Plugin defines the interface that all data extraction plugins must implement
type Plugin interface {
	// Extract performs the data extraction and returns a plugin response
	Extract(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error)
	ID() string
}

// NewPlugin creates a new plugin instance from config
func NewPlugin(config *PluginConfig) (Plugin, error) {
	return &binaryPlugin{
		config: config,
	}, nil
}

// binaryPlugin implements Plugin interface for binary plugins
type binaryPlugin struct {
	config *PluginConfig
}

// Extract implements Plugin interface
func (p *binaryPlugin) Extract(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	// Set plugin ID in request
	req.PluginId = p.config.ID

	// For now, just return empty response
	return &pb.PluginResponse{
		PluginId: p.config.ID,
		Items:    []*pb.DataItem{},
	}, nil
}

func (p *binaryPlugin) ID() string {
	return p.config.ID
}
