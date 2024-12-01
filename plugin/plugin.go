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

// NewPlugin creates a new plugin instance from config
func NewPlugin(config *PluginConfig) (Plugin, error) {
	return &binaryPlugin{
		config: config,
	}, nil
}

func (p *binaryPlugin) ID() string {
	return p.config.ID
}

func (p *binaryPlugin) Configure(config map[string]string) error {
	p.config.Config = config
	return nil
}

// Extract implements Plugin interface
func (p *binaryPlugin) Extract(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	// Set plugin ID and operation in request
	req.PluginId = p.config.ID
	req.Operation = "extract"

	// Serialize request using protobuf
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Start plugin process
	cmd := exec.CommandContext(ctx, p.config.BinaryPath)
	cmd.Stdin = bytes.NewReader(reqData)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run plugin
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run plugin: %v\nstderr: %s", err, stderr.String())
	}

	// Parse response using protobuf
	resp := &pb.PluginResponse{}
	if err := proto.Unmarshal(stdout.Bytes(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return resp, nil
}
