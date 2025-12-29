package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the execman configuration.
type Config struct {
	DefaultInstallDir  string `json:"default_install_dir,omitempty"`
	IncludePrereleases bool   `json:"include_prereleases"`
	path               string // internal, not serialized
}

// DefaultConfigPath returns the default config file path.
func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(configDir, "execman", "config.json"), nil
}

// Load loads the config from the default location.
func Load() (*Config, error) {
	path, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

// LoadFrom loads the config from a specific path.
func LoadFrom(path string) (*Config, error) {
	// If file doesn't exist, return a new config with defaults.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		return &Config{
			DefaultInstallDir:  filepath.Join(homeDir, ".local", "bin"),
			IncludePrereleases: false,
			path:               path,
		}, nil
	}

	// #nosec G304 -- Reading user config from trusted path
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.path = path

	// Set defaults if not specified.
	if cfg.DefaultInstallDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		cfg.DefaultInstallDir = filepath.Join(homeDir, ".local", "bin")
	}

	return &cfg, nil
}

// Save saves the config to disk.
func (c *Config) Save() error {
	// Ensure directory exists.
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Use 0600 permissions for config file (user read/write only).
	if err := os.WriteFile(c.path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
