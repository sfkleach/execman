package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Executable represents a managed executable in the registry
type Executable struct {
	Source      string    `json:"source"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
	Path        string    `json:"path"`
	Platform    string    `json:"platform"`
	Checksum    string    `json:"checksum"`
}

// Registry represents the execman registry
type Registry struct {
	SchemaVersion     int                    `json:"schema_version"`
	DefaultInstallDir string                 `json:"default_install_dir"`
	Executables       map[string]*Executable `json:"executables"`
	path              string                 // internal, not serialized
}

// DefaultRegistryPath returns the default registry file path
func DefaultRegistryPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(configDir, "execman", "registry.json"), nil
}

// Load loads the registry from the default location
func Load() (*Registry, error) {
	path, err := DefaultRegistryPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

// LoadFrom loads the registry from a specific path
func LoadFrom(path string) (*Registry, error) {
	// If file doesn't exist, return a new empty registry
	if _, err := os.Stat(path); os.IsNotExist(err) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		defaultInstallDir := filepath.Join(homeDir, ".local", "bin")
		return &Registry{
			SchemaVersion:     1,
			DefaultInstallDir: defaultInstallDir,
			Executables:       make(map[string]*Executable),
			path:              path,
		}, nil
	}

	// #nosec G304 -- Reading user registry from trusted path
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	reg.path = path
	if reg.Executables == nil {
		reg.Executables = make(map[string]*Executable)
	}

	return &reg, nil
}

// Save saves the registry to disk
func (r *Registry) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	// Use 0600 permissions for registry file (user read/write only)
	if err := os.WriteFile(r.path, data, 0600); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// Add adds or updates an executable in the registry
func (r *Registry) Add(name string, exec *Executable) {
	r.Executables[name] = exec
}

// Get retrieves an executable from the registry
func (r *Registry) Get(name string) (*Executable, bool) {
	exec, ok := r.Executables[name]
	return exec, ok
}

// Remove removes an executable from the registry
func (r *Registry) Remove(name string) {
	delete(r.Executables, name)
}

// List returns all executable names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.Executables))
	for name := range r.Executables {
		names = append(names, name)
	}
	return names
}
