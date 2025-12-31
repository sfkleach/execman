package list

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sfkleach/execman/pkg/registry"
	"github.com/spf13/cobra"
)

// ListOutput represents the JSON output format for the list command.
type ListOutput struct {
	Executables []ExecutableInfo `json:"executables"`
}

// ExecutableInfo represents information about a single executable.
type ExecutableInfo struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Version     string `json:"version"`
	Path        string `json:"path"`
	InstalledAt string `json:"installed_at"`
}

// NewListCommand creates the list command.
func NewListCommand() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all managed executables",
		Long:  "Display all executables managed by execman with their versions and locations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}

func runList(jsonOutput bool) error {
	// Load registry.
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Get all executable names.
	names := reg.List()

	if len(names) == 0 {
		if jsonOutput {
			output := ListOutput{Executables: []ExecutableInfo{}}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(output)
		}
		fmt.Println("No managed executables.")
		return nil
	}

	// Sort by name.
	sort.Strings(names)

	if jsonOutput {
		return outputJSON(reg, names)
	}

	return outputText(reg, names)
}

func outputJSON(reg *registry.Registry, names []string) error {
	// Convert to output format.
	executables := make([]ExecutableInfo, 0, len(names))
	for _, name := range names {
		exec, ok := reg.Get(name)
		if !ok {
			continue
		}

		executables = append(executables, ExecutableInfo{
			Name:        name,
			Source:      exec.Source,
			Version:     exec.Version,
			Path:        exec.Path,
			InstalledAt: exec.InstalledAt.Format(time.RFC3339),
		})
	}

	output := ListOutput{Executables: executables}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func outputText(reg *registry.Registry, names []string) error {
	fmt.Println("Managed executables:")
	fmt.Println()

	homeDir, _ := os.UserHomeDir()

	for _, name := range names {
		exec, ok := reg.Get(name)
		if !ok {
			continue
		}

		// Display path with ~ for home directory.
		displayPath := exec.Path
		if homeDir != "" && strings.HasPrefix(exec.Path, homeDir) {
			displayPath = "~" + strings.TrimPrefix(exec.Path, homeDir)
		}

		// Extract repo path from source URL.
		source := strings.TrimPrefix(exec.Source, "https://")

		// Format installed_at timestamp.
		installedDate := exec.InstalledAt.Format("2006-01-02")

		// Get just the executable name from the path.
		execName := filepath.Base(exec.Path)

		// Print formatted output.
		fmt.Printf("  %-15s %-9s %s\n", execName, exec.Version, displayPath)
		fmt.Printf("  %-15s %-9s %s\n", "", "", source)
		fmt.Printf("  %-15s %-9s installed %s\n", "", "", installedDate)
		fmt.Println()
	}

	count := len(names)
	if count == 1 {
		fmt.Println("1 executable managed")
	} else {
		fmt.Printf("%d executables managed\n", count)
	}

	return nil
}
