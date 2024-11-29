package plugin

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/encoding/protojson"
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
	// Set plugin ID and operation in request
	req.PluginId = p.config.ID
	req.Operation = "extract"

	// Serialize request to JSON
	reqJson, err := protojson.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, p.config.BinaryPath)
	cmd.Stdin = bytes.NewReader(reqJson)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run plugin
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run plugin: %v\nstderr: %s", err, stderr.String())
	}

	// Parse response
	resp := &pb.PluginResponse{}
	if err := protojson.Unmarshal(stdout.Bytes(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return resp, nil

}

func (p *binaryPlugin) ID() string {
	return p.config.ID
}
