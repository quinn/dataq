package plugin

import (
	"context"

	pb "go.quinn.io/dataq/proto"
)

// PluginConfig contains configuration for a plugin
type PluginConfig struct {
	ID         string            `yaml:"id"`
	Name       string            `yaml:"name"`
	BinaryPath string            `yaml:"binary_path"`
	Config     map[string]string `yaml:"config"`
	Enabled    bool              `yaml:"enabled"`
}

// Plugin defines the interface that all data extraction plugins must implement
type Plugin interface {
	ID() string
	Configure(map[string]string) error
	Extract(context.Context, *pb.PluginRequest) (*pb.PluginResponse, error)
}

// binaryPlugin implements Plugin using a binary executable
type binaryPlugin struct {
	config *PluginConfig
}

func (p *binaryPlugin) ID() string {
	return p.config.ID
}

func (p *binaryPlugin) Configure(config map[string]string) error {
	p.config.Config = config
	return nil
}
