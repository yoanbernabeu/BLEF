#!/bin/bash

# BLEF CLI Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/yoanbernabeu/BLEF/main/tools/blef-cli/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="yoanbernabeu/BLEF"
BINARY_NAME="blef-cli"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

echo ""
echo -e "${BLUE}=========================================="
echo "  BLEF CLI Installer"
echo -e "==========================================${NC}"
echo ""

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            echo -e "${RED}❌ Error: Unsupported operating system: $OS${NC}"
            echo ""
            echo "This installer supports Linux and macOS only."
            echo "Windows users: Download manually from https://github.com/$REPO/releases"
            echo ""
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}❌ Error: Unsupported architecture: $ARCH${NC}"
            echo ""
            exit 1
            ;;
    esac

    echo -e "${GREEN}✓ Detected platform: $OS-$ARCH${NC}"
}

# Get latest release version
get_latest_version() {
    echo "🔍 Fetching latest version..."
    
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        echo -e "${RED}❌ Error: Could not fetch latest version${NC}"
        echo "Please check your internet connection or try again later."
        exit 1
    fi
    
    echo -e "${GREEN}✓ Latest version: $VERSION${NC}"
}

# Download and install
install_binary() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/blef-cli-$VERSION-$OS-$ARCH.tar.gz"
    
    echo "📥 Downloading $BINARY_NAME $VERSION..."
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    # Download with progress
    if ! curl -fsSL --progress-bar "$DOWNLOAD_URL" -o "$TMP_DIR/blef-cli.tar.gz"; then
        echo -e "${RED}❌ Error: Failed to download $BINARY_NAME${NC}"
        echo "URL: $DOWNLOAD_URL"
        echo ""
        echo "Please check:"
        echo "  1. Your internet connection"
        echo "  2. The release exists: https://github.com/$REPO/releases"
        echo ""
        exit 1
    fi
    
    # Extract
    echo "📦 Extracting..."
    tar -xzf "$TMP_DIR/blef-cli.tar.gz" -C "$TMP_DIR"
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Install binary
    echo "⚙️  Installing to $INSTALL_DIR..."
    mv "$TMP_DIR/blef-cli-$OS-$ARCH" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    echo -e "${GREEN}✓ $BINARY_NAME installed successfully!${NC}"
}

# Check if directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo ""
        echo -e "${YELLOW}⚠️  Warning: $INSTALL_DIR is not in your PATH${NC}"
        echo ""
        echo "To use blef-cli from anywhere, add this line to your shell config:"
        echo ""
        
        # Detect shell
        if [ -n "$ZSH_VERSION" ]; then
            echo "    echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc"
            echo "    source ~/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            echo "    echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc"
            echo "    source ~/.bashrc"
        else
            echo "    export PATH=\"$INSTALL_DIR:\$PATH\""
        fi
        echo ""
        
        return 1
    fi
    return 0
}

# Verify installation
verify_installation() {
    echo ""
    echo "🔍 Verifying installation..."
    
    if command -v $BINARY_NAME &> /dev/null; then
        VERSION_OUTPUT=$($BINARY_NAME --version 2>&1 || echo "unknown")
        echo -e "${GREEN}✓ Installation verified!${NC}"
        echo "Version: $VERSION_OUTPUT"
        PATH_IN_PATH=true
    elif [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        VERSION_OUTPUT=$("$INSTALL_DIR/$BINARY_NAME" --version 2>&1 || echo "unknown")
        echo -e "${GREEN}✓ Binary installed!${NC}"
        echo "Version: $VERSION_OUTPUT"
        check_path
        PATH_IN_PATH=false
    else
        echo -e "${RED}❌ Installation verification failed${NC}"
        exit 1
    fi
}

# Show usage examples
show_examples() {
    echo ""
    echo -e "${BLUE}=========================================="
    echo "  Quick Start"
    echo -e "==========================================${NC}"
    echo ""
    echo "Validate a BLEF file:"
    echo "  $ blef-cli validate my-library.blef.json"
    echo ""
    echo "Convert a Goodreads export:"
    echo "  $ blef-cli convert goodreads_export.csv"
    echo ""
    echo "View your library interactively:"
    echo "  $ blef-cli view my-library.blef.json"
    echo ""
    echo "Get help:"
    echo "  $ blef-cli --help"
    echo ""
}

# Show links
show_links() {
    echo -e "${BLUE}📚 Documentation:${NC}"
    echo "  https://github.com/$REPO"
    echo ""
    echo -e "${BLUE}🐛 Issues:${NC}"
    echo "  https://github.com/$REPO/issues"
    echo ""
}

# Main installation flow
main() {
    detect_platform
    get_latest_version
    install_binary
    verify_installation
    show_examples
    show_links
    
    echo -e "${GREEN}✨ Installation complete!${NC}"
    echo ""
}

main
