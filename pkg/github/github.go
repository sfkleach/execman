package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Release represents a GitHub release
type Release struct {
	TagName    string  `json:"tag_name"`
	Name       string  `json:"name"`
	Prerelease bool    `json:"prerelease"`
	Assets     []Asset `json:"assets"`
}

// Asset represents a GitHub release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// ParseSource parses a GitHub source string into owner, repo, and optional version
func ParseSource(source string) (owner, repo, version string, err error) {
	// Support formats:
	// - github.com/owner/repo
	// - github.com/owner/repo@version
	// - https://github.com/owner/repo
	// - https://github.com/owner/repo@version

	source = strings.TrimPrefix(source, "https://")
	source = strings.TrimPrefix(source, "http://")
	source = strings.TrimPrefix(source, "github.com/")

	// Check for version suffix
	if strings.Contains(source, "@") {
		parts := strings.SplitN(source, "@", 2)
		source = parts[0]
		version = parts[1]
	}

	// Parse owner/repo
	parts := strings.Split(source, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid GitHub source format: %s", source)
	}

	owner = parts[0]
	repo = parts[1]

	return owner, repo, version, nil
}

// ToURL converts owner/repo to a GitHub URL
func ToURL(owner, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s", owner, repo)
}

// GetLatestRelease fetches the latest release from GitHub
func GetLatestRelease(owner, repo string, includePrereleases bool) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	// #nosec G107 -- URL is constructed from validated GitHub repo components
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to parse releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", owner, repo)
	}

	// Find the first non-prerelease (or first release if includePrereleases)
	for _, release := range releases {
		if !release.Prerelease || includePrereleases {
			return &release, nil
		}
	}

	return nil, fmt.Errorf("no suitable releases found for %s/%s", owner, repo)
}

// GetRelease fetches a specific release by tag from GitHub
func GetRelease(owner, repo, tag string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)

	// #nosec G107 -- URL is constructed from validated GitHub repo components
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return &release, nil
}

// FindAsset finds a matching asset for the given OS and architecture
func FindAsset(assets []Asset, osName, arch string) (*Asset, error) {
	// Pattern variations to match
	patterns := []string{
		fmt.Sprintf("_%s_%s\\.tar\\.gz$", osName, arch),
		fmt.Sprintf("-%s-%s\\.tar\\.gz$", osName, arch),
		fmt.Sprintf("_%s_%s$", osName, arch),
		fmt.Sprintf("-%s-%s$", osName, arch),
	}

	for _, asset := range assets {
		assetLower := strings.ToLower(asset.Name)
		for _, pattern := range patterns {
			matched, _ := regexp.MatchString(pattern, assetLower)
			if matched {
				return &asset, nil
			}
		}
	}

	return nil, fmt.Errorf("no matching asset found for %s/%s", osName, arch)
}

// DownloadAsset downloads an asset from GitHub
func DownloadAsset(asset *Asset, dest string) error {
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read asset data: %w", err)
	}

	// Use 0600 permissions for downloaded file (temp file)
	if err := os.WriteFile(dest, out, 0600); err != nil {
		return fmt.Errorf("failed to write asset to file: %w", err)
	}

	return nil
}
