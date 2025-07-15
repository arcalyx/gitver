#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Gitver Installation Script${NC}"
echo "==============================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Create build directory if it doesn't exist
mkdir -p build

# Build the binary
echo "Building gitver..."
go build -o build/gitver
if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
echo -e "${GREEN}Build successful!${NC}"

# Determine installation path based on platform
INSTALL_DIR=""
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    INSTALL_DIR="/usr/local/bin"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    INSTALL_DIR="/usr/local/bin"
elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo -e "${YELLOW}Windows detected. Installing to ./build only.${NC}"
    echo -e "${GREEN}Installation successful! Binary located at ./build/gitver${NC}"
    echo "Add this directory to your PATH to use gitver from anywhere."
    exit 0
else
    echo -e "${YELLOW}Unknown OS. Installing to ./build only.${NC}"
    echo -e "${GREEN}Installation successful! Binary located at ./build/gitver${NC}"
    echo "Add this directory to your PATH to use gitver from anywhere."
    exit 0
fi

# Install the binary
echo "Installing gitver to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    # User has write permissions to the install directory
    cp build/gitver "$INSTALL_DIR/gitver"
else
    # Need sudo to install
    echo "Requesting administrator privileges to install to $INSTALL_DIR"
    sudo cp build/gitver "$INSTALL_DIR/gitver"
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}Installation failed!${NC}"
    echo "You can manually copy the binary from ./build/gitver to a directory in your PATH."
    exit 1
fi

echo -e "${GREEN}Installation successful!${NC}"
echo "You can now use gitver from anywhere."
echo "Run 'gitver --help' to get started."
