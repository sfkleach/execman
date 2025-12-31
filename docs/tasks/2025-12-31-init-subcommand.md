# init subcommand

**Goal**: A streamlined installation of `execman` or both `execman` and `pathman`.

## The subcommand

- `execman init FOLDER` will create the config.json and registry.json file, setting the default folder to be FOLDER, then run the equivalent of `execman install github.com/sfkleach/execman`. This will ensure that the current version of execman is installed in the expected location with the appropriate meta-info.

## The install.sh and install-with-pathman.sh scripts

- The install.sh script is used to:
  1. download execman from https://github.com/sfkleach/execman/releases/latest/download/execman_linux_amd64.tar.gz (with platform detection)
  2. run `./execman init ~/.local/bin`
  3. `rm ./execman` (unless the current folder is `~/.local/bin`)

- The install-with-pathman.sh script will:
  1. download execman from https://github.com/sfkleach/execman/releases/latest/download/execman_linux_amd64.tar.gz (with platform detection)
  2. run  `./execman init ~/.local/bin`
  3. use `~/.local/bin/execman` to install `pathman`
  4. use `~/.local/bin/pathman init` to patch the $PATH script into .profile and create configuration files
  5. use `pathman add ~/.local/bin` to ensure they are both on the $PATH
  6. `rm ./execman` (unless the current folder is `~/.local/bin`)



## Note on pathman's init subcommand

```
pathman on  update-subcommand [$!] via 🐹 v1.24.2 took 4s 
❯ go run ./cmd/pathman init --help
Create the managed folder with appropriate permissions.
If the folder already exists, check its permissions and warn if insecure.

Use --no for non-interactive mode (suitable for scripts). In non-interactive
mode, only the folder structure is created - no shell profile modifications
or binary relocations are performed.

Usage:
  pathman init [flags]

Flags:
  -h, --help   help for init
      --no     Non-interactive mode: create folders only, no prompts
```