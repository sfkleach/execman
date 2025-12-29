# Execman

Execman is a command-line tool for managing standalone executables from GitHub releases.

## Features

- **Install** executables directly from GitHub releases
- **Track** installed executables with version and origin information
- **Registry** maintains metadata for secure updates
- **Cross-platform** support for Linux, macOS, and Windows

## Building

```bash
go build -o bin/execman ./cmd/execman
```

Or using the Justfile:

```bash
just build
```

## Usage

### Install an executable

```bash
# Install latest version
execman install github.com/owner/repo

# Install specific version
execman install github.com/owner/repo@v1.2.3

# Install to custom directory
execman install github.com/owner/repo --into /usr/local/bin

# Skip confirmation prompts
execman install github.com/owner/repo --yes
```

### Show version

```bash
# Show version using flag
./execman --version

# Show version using subcommand
./execman version
```

### Get help

```bash
# Show all commands
execman

# Get help for any command
execman [command] --help
```

## Commands

- `version` - Print the version number of execman
- `install` - Install an executable from GitHub releases
- `list` - List managed executables (TBD)
- `info` - Show information about an executable (TBD)
- `check` - Check executable status (TBD)
- `update` - Update an executable (TBD)
- `remove` - Remove an executable (TBD)
- `adopt` - Adopt an existing executable (TBD)

## Configuration

### Registry

Location: `~/.config/execman/registry.json`

Tracks all installed executables with version, source, checksum, and path information.

### Config (Optional)

Location: `~/.config/execman/config.json`

```json
{
  "default_install_dir": "/home/user/.local/bin",
  "include_prereleases": false
}
```

Defaults:
- `default_install_dir`: `~/.local/bin`
- `include_prereleases`: `false`

## Project Structure

```
execman/
├── cmd/
│   └── execman/
│       └── main.go          # Main entry point
├── pkg/
│   ├── archive/             # Archive extraction and checksums
│   ├── config/              # Configuration management
│   ├── github/              # GitHub API integration
│   ├── install/             # Install command implementation
│   ├── registry/            # Registry management
│   └── version/             # Version information
├── go.mod
└── README.md
```
