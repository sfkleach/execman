# Execman

Execman is a command-line tool for managing executables.

## Building

```bash
go build -o execman ./cmd/execman
```

## Usage

```bash
# Show help
./execman

# Show version using flag
./execman --version

# Show version using subcommand
./execman version

# Get help for any command
./execman [command] --help
```

## Commands

- `version` - Print the version number of execman
- `install` - Install an executable (TBD)
- `info` - Show information about an executable (TBD)
- `list` - List managed executables (TBD)
- `check` - Check executable status (TBD)
- `update` - Update an executable (TBD)
- `remove` - Remove an executable (TBD)
- `adopt` - Adopt an existing executable (TBD)

## Project Structure

```
execman/
├── cmd/
│   └── execman/
│       └── main.go          # Main entry point
├── pkg/
│   └── version/
│       └── version.go       # Version information
├── go.mod
└── README.md
```
