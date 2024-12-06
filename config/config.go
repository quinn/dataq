package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	dirname = "dataq"
	config  = "config.yaml"
)

type Config struct {
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

func ConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get home directory: %v", err)
		}

		configDir = filepath.Join(home, ".config")
	}

	return filepath.Join(configDir, dirname)
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), config)
}

func DataDir() string {
	return filepath.Join(ConfigDir(), "data")
}

func StateDir() string {
	return filepath.Join(ConfigDir(), "state")
}

func Get() (*Config, error) {
	configData, err := os.ReadFile(ConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
