package github

import (
	"testing"
)

func TestFindAsset(t *testing.T) {
	tests := []struct {
		name      string
		assets    []Asset
		osName    string
		arch      string
		wantName  string
		wantError bool
	}{
		{
			name: "Linux x86_64 tar.gz with underscore separator",
			assets: []Asset{
				{Name: "nutmeg-compiler_Linux_x86_64.tar.gz"},
				{Name: "nutmeg-compiler_Darwin_x86_64.tar.gz"},
			},
			osName:   "linux",
			arch:     "amd64",
			wantName: "nutmeg-compiler_Linux_x86_64.tar.gz",
		},
		{
			name: "Linux amd64 tar.gz with hyphen separator",
			assets: []Asset{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
			},
			osName:   "linux",
			arch:     "amd64",
			wantName: "tool-linux-amd64.tar.gz",
		},
		{
			name: "Darwin arm64 zip",
			assets: []Asset{
				{Name: "app_Darwin_arm64.zip"},
				{Name: "app_Linux_arm64.zip"},
			},
			osName:   "darwin",
			arch:     "arm64",
			wantName: "app_Darwin_arm64.zip",
		},
		{
			name: "Windows amd64 zip",
			assets: []Asset{
				{Name: "tool_Windows_x86_64.zip"},
				{Name: "tool_Linux_x86_64.tar.gz"},
			},
			osName:   "windows",
			arch:     "amd64",
			wantName: "tool_Windows_x86_64.zip",
		},
		{
			name: "Linux arm64 with aarch64 alias",
			assets: []Asset{
				{Name: "binary_linux_aarch64.tar.gz"},
				{Name: "binary_linux_x86_64.tar.gz"},
			},
			osName:   "linux",
			arch:     "arm64",
			wantName: "binary_linux_aarch64.tar.gz",
		},
		{
			name: "Case insensitive OS matching",
			assets: []Asset{
				{Name: "app_LINUX_AMD64.tar.gz"},
			},
			osName:   "linux",
			arch:     "amd64",
			wantName: "app_LINUX_AMD64.tar.gz",
		},
		{
			name: "No extension",
			assets: []Asset{
				{Name: "binary_linux_amd64"},
			},
			osName:   "linux",
			arch:     "amd64",
			wantName: "binary_linux_amd64",
		},
		{
			name: "386 architecture with i386 alias",
			assets: []Asset{
				{Name: "tool_linux_i386.tar.gz"},
			},
			osName:   "linux",
			arch:     "386",
			wantName: "tool_linux_i386.tar.gz",
		},
		{
			name: "No matching asset",
			assets: []Asset{
				{Name: "tool_darwin_amd64.tar.gz"},
				{Name: "tool_windows_amd64.zip"},
			},
			osName:    "linux",
			arch:      "amd64",
			wantError: true,
		},
		{
			name: "Skip checksums file",
			assets: []Asset{
				{Name: "checksums.txt"},
				{Name: "app_linux_amd64.tar.gz"},
			},
			osName:   "linux",
			arch:     "amd64",
			wantName: "app_linux_amd64.tar.gz",
		},
		{
			name:      "Empty assets list",
			assets:    []Asset{},
			osName:    "linux",
			arch:      "amd64",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := FindAsset(tt.assets, tt.osName, tt.arch)

			if tt.wantError {
				if err == nil {
					t.Errorf("FindAsset() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("FindAsset() unexpected error: %v", err)
				return
			}

			if asset.Name != tt.wantName {
				t.Errorf("FindAsset() = %q, want %q", asset.Name, tt.wantName)
			}
		})
	}
}

func TestParseSource(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		wantOwner   string
		wantRepo    string
		wantVersion string
		wantError   bool
	}{
		{
			name:      "Simple owner/repo",
			source:    "github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:        "With version",
			source:      "github.com/owner/repo@v1.2.3",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "v1.2.3",
		},
		{
			name:      "With https prefix",
			source:    "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:        "With https and version",
			source:      "https://github.com/owner/repo@v1.0.0",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "v1.0.0",
		},
		{
			name:      "Without github.com prefix",
			source:    "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "With http prefix",
			source:    "http://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "Invalid format - no slash",
			source:    "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, version, err := ParseSource(tt.source)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseSource() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSource() unexpected error: %v", err)
				return
			}

			if owner != tt.wantOwner {
				t.Errorf("ParseSource() owner = %q, want %q", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("ParseSource() repo = %q, want %q", repo, tt.wantRepo)
			}
			if version != tt.wantVersion {
				t.Errorf("ParseSource() version = %q, want %q", version, tt.wantVersion)
			}
		})
	}
}
