#!/bin/bash

# SaveSync Installation Script

set -e

echo "======================================"
echo "  SaveSync Installation"
echo "======================================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${BLUE}Detected Go version:${NC}"
go version
echo ""

# Get GOPATH
GOPATH=$(go env GOPATH)
if [ -z "$GOPATH" ]; then
    echo -e "${RED}Error: GOPATH not set${NC}"
    exit 1
fi

INSTALL_DIR="$GOPATH/bin"

echo -e "${BLUE}Installation directory: ${NC}$INSTALL_DIR"
echo ""

# Check if directory exists
if [ ! -d "$INSTALL_DIR" ]; then
    echo -e "${BLUE}Creating installation directory...${NC}"
    mkdir -p "$INSTALL_DIR"
fi

# Build the binary
echo -e "${BLUE}Building SaveSync...${NC}"
go build -o savesync cmd/main.go

if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Install binary
echo -e "${BLUE}Installing to $INSTALL_DIR...${NC}"
mv savesync "$INSTALL_DIR/savesync"
chmod +x "$INSTALL_DIR/savesync"

echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""

# Check if GOPATH/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${RED}Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add the following line to your shell configuration file:"
    echo "  (~/.bashrc, ~/.zshrc, etc.)"
    echo ""
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
    echo "Then run: source ~/.bashrc (or your shell config file)"
    echo ""
else
    echo -e "${GREEN}✓ $INSTALL_DIR is in PATH${NC}"
    echo ""
fi

# Test installation
echo -e "${BLUE}Testing installation...${NC}"
"$INSTALL_DIR/savesync" version

echo ""
echo -e "${GREEN}======================================"
echo "  Installation Complete!"
echo "======================================${NC}"
echo ""
echo "Get started with:"
echo "  savesync help"
echo ""
echo "Example usage:"
echo "  savesync add-game --name \"My Game\" --path \"/path/to/saves\""
echo "  savesync checkpoint --game mygame --name \"Important Save\""
echo ""
