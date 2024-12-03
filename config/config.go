package config

type Config struct {
	DataDir string    `yaml:"data_dir"`
	Plugins []*Plugin `yaml:"plugins"`
}

// PluginConfig contains configuration for a plugin
type Plugin struct {
	ID         string            `yaml:"id"`
	Name       string            `yaml:"name"`
	BinaryPath string            `yaml:"binary_path"`
	Config     map[string]string `yaml:"config"`
	Enabled    bool              `yaml:"enabled"`
}
