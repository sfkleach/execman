# Execman

Execman is a command-line tool for managing standalone executables from GitHub releases. It can be used on its own or paired with the $PATH management utility `pathman`, which can be installed at the same time.

## Quick Start

### First Time Installation

**Option 1: Using the install script** (recommended):

```bash
curl -sSL https://raw.githubusercontent.com/sfkleach/execman/main/scripts/install.sh | bash
```

This downloads execman and initializes it in `~/.local/bin`.

**Option 2: Using Go**:

```bash
go install github.com/sfkleach/execman/cmd/execman@latest
```

**Option 3: Download from releases**:

Download the appropriate archive for your platform from the [releases page](https://github.com/sfkleach/execman/releases/latest), extract it, and run:

```bash
./execman init ~/.local/bin
```

### With Pathman Integration

To also install [pathman](https://github.com/sfkleach/pathman) and automatically configure your PATH:

```bash
curl -sSL https://raw.githubusercontent.com/sfkleach/execman/main/scripts/install-with-pathman.sh | bash
```

## Features

- **Install** executables directly from GitHub releases
- **Track** installed executables with version and origin information
- **List** all managed executables with details
- **Check** for available updates across all executables
- **Update** executables individually or all at once
- **Remove** executables and delete files
- **Forget** executables while keeping files on disk
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

### Initialize execman

```bash
# Initialize configuration and install execman itself
execman init <folder>

# Example: Initialize in ~/.local/bin
execman init ~/.local/bin
```

This creates the configuration and registry files, then installs execman itself from GitHub with proper metadata.

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

### List managed executables

```bash
# Show all managed executables
execman list
execman ls

# Show specific executable
execman list myapp

# Show detailed information
execman list --long
execman ls -l

# Show specific executable with details
execman list myapp --long

# Output as JSON
execman list --json
```

### Check for updates

```bash
# Check all executables for updates and integrity
execman check

# Check specific executable
execman check myapp

# Show all executables including up-to-date ones
execman check --no-skip

# Verify checksums of installed executables
execman check --verify

# Output as JSON
execman check --json
```

### Update executables

```bash
# Update a specific executable
execman update myapp

# Update all executables
execman update --all

# Skip confirmation prompts
execman update --all --yes

# Reinstall a missing executable
execman update myapp  # Will detect missing file and offer reinstall
```

### Remove an executable

```bash
# Remove executable and delete file
execman remove myapp

# Skip confirmation prompt
execman remove myapp --yes
```

### Forget an executable

```bash
# Stop tracking but keep the file
execman forget myapp

# Skip confirmation prompt
execman forget myapp --yes
```

### Show version

```bash
# Show version using flag
execman --version

# Show version using subcommand
execman version
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
- `init` - Initialize execman configuration and install execman itself
- `install` - Install an executable from GitHub releases
- `list` (alias: `ls`) - List managed executables with optional filtering and detailed view
- `check` - Check for available updates and verify integrity
- `update` - Update executables to latest versions
- `remove` - Remove an executable and delete the file
- `forget` - Stop tracking an executable but keep the file

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

## Example Workflow

```bash
# Install some executables
execman install github.com/sfkleach/pathman --yes
execman install github.com/sfkleach/nutmeg-run --yes

# List what's installed
execman list

# Check for updates
execman check

# Update all executables
execman update --all --yes

# Remove an executable
execman remove nutmeg-run

# Stop tracking but keep the file
execman forget pathman
```

## Uninstalling

To remove execman from your system, see the [uninstall instructions](docs/uninstall-execman.md).

## Project Structure

```
execman/
├── cmd/
│   └── execman/
│       └── main.go          # Main entry point
├── pkg/
│   ├── archive/             # Archive extraction and checksums
│   ├── check/               # Check command implementation
│   ├── config/              # Configuration management
│   ├── forget/              # Forget command implementation
│   ├── github/              # GitHub API integration
│   ├── init/                # Init command implementation
│   ├── install/             # Install command implementation
│   ├── list/                # List command implementation
│   ├── registry/            # Registry management
│   ├── remove/              # Remove command implementation
│   ├── symlink/             # Symlink detection and handling
│   ├── update/              # Update command implementation
│   └── version/             # Version information
├── scripts/
│   ├── install.sh           # Installation script
│   └── install-with-pathman.sh  # Installation script with pathman
├── go.mod
└── README.md
```

## Symlink Handling

When updating or removing an executable that is a symbolic link, execman will prompt you to choose how to handle it:

```
Note: /usr/local/bin/myapp is a symlink to /opt/myapp/v1.2.3/myapp

How would you like to proceed?
  [1] Replace the symlink target (/opt/myapp/v1.2.3/myapp)
  [2] Replace the symlink itself (/usr/local/bin/myapp)
  [3] Cancel

Choice [1/2/3]:
```

- **Option 1**: Operates on the target file that the symlink points to
- **Option 2**: Removes the symlink and operates on that location directly
- **Option 3**: Cancels the operation

In non-interactive mode (`--yes`), symlink operations will fail with an error message instructing you to run without `--yes` to choose how to handle the symlink.
