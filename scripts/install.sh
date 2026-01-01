#!/bin/bash
# Install script for execman
# Downloads execman and initializes it in ~/.local/bin

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

# Clean up the bootstrap execman unless we're already in the install directory.
CURRENT_DIR=$(pwd)
if [ "$CURRENT_DIR" != "$INSTALL_DIR" ]; then
    rm -f "$TEMP_DIR/execman"
fi

echo ""
echo "Installation complete!"
