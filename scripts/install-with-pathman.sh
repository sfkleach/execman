#!/bin/bash
# Install script for execman with pathman integration
# Downloads execman, initializes it, installs pathman, and sets up PATH

set -e

INSTALL_DIR="$HOME/.local/bin"
REPO="sfkleach/execman"
TEMP_DIR=$(mktemp -d)

cleanup() {
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Detect OS and architecture.
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    i386|i686)
        ARCH="386"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Adjust OS name for darwin.
if [ "$OS" = "darwin" ]; then
    OS="darwin"
fi

echo "Downloading execman for $OS/$ARCH..."
URL="https://github.com/$REPO/releases/latest/download/execman_${OS}_${ARCH}.tar.gz"

if ! curl -L -f -o "$TEMP_DIR/execman.tar.gz" "$URL"; then
    echo "Error: Failed to download execman from $URL"
    exit 1
fi

echo "Extracting..."
tar xzf "$TEMP_DIR/execman.tar.gz" -C "$TEMP_DIR"

echo ""
echo "Running execman init..."
"$TEMP_DIR/execman" init "$INSTALL_DIR"

echo ""
echo "Installing pathman..."
"$INSTALL_DIR/execman" install github.com/sfkleach/pathman --yes

echo ""
echo "Setting up PATH with pathman..."
"$INSTALL_DIR/pathman" init --no

echo ""
echo "Adding $INSTALL_DIR to PATH..."
"$INSTALL_DIR/pathman" add "$INSTALL_DIR"

# Clean up the bootstrap execman unless we're already in the install directory.
CURRENT_DIR=$(pwd)
if [ "$CURRENT_DIR" != "$INSTALL_DIR" ]; then
    rm -f "$TEMP_DIR/execman"
fi

echo ""
echo "Installation complete!"
echo ""
echo "Both execman and pathman are now installed and $INSTALL_DIR is on your PATH."
echo "You may need to restart your shell or source your profile for PATH changes to take effect."
