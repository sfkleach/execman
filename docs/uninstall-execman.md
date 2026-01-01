# Uninstall execman

To completely remove `execman` from your system, follow these instructions:

## 1. Remove All Managed Executables (Optional)

If you want to also remove all executables that execman manages:

```bash
# List all managed executables
execman list

# Remove each one (this deletes the files)
execman remove <executable-name> --yes

# Or just forget them (keeps the files)
execman forget <executable-name> --yes
```

## 2. Remove execman Itself

Remove the execman binary from your system:

```bash
# Find where execman is installed
which execman
# or
execman list execman

# Remove it (typical location is ~/.local/bin/execman)
rm ~/.local/bin/execman
# or wherever your installation is located
```

## 3. Remove Configuration and Registry

Remove execman's configuration and registry files:

```bash
rm -rf ~/.config/execman
```

This removes:
- `~/.config/execman/config.json` - configuration file
- `~/.config/execman/registry.json` - registry of managed executables

## 4. Clean Up PATH (If Using Pathman)

If you used the `install-with-pathman.sh` script, you may want to:

```bash
# Remove ~/.local/bin from PATH using pathman
pathman remove ~/.local/bin

# Or manually edit your shell profile
# Remove the line that pathman added from ~/.bashrc or ~/.profile
```

## Complete One-Liner (Standard Installation)

For a standard installation in `~/.local/bin`:

```bash
rm ~/.local/bin/execman && rm -rf ~/.config/execman
```

## Verification

Verify execman is completely removed:

```bash
# Should return nothing or "not found"
which execman

# Should not exist
ls ~/.config/execman
```
