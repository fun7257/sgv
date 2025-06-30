#!/bin/bash
#
# sgv (Simple Go Version) - Installation Script
#
# This script downloads and installs the latest version of sgv for your system.
# It is designed to work on macOS and Linux systems.
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/fun7257/sgv/main/install.sh | bash
#

set -e # Exit immediately if a command exits with a non-zero status.

# --- Helper Functions ---

# Color codes for messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

# --- Main Installation Logic ---

main() {
    # 1. Detect OS and Architecture
    OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH_TYPE=$(uname -m)
    INSTALL_DIR="/usr/local/bin"
    
    case "$OS_TYPE" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            error "Unsupported OS: $OS_TYPE. Only Linux and macOS are supported."
            ;;
    esac

    case "$ARCH_TYPE" in
        x86_64 | amd64)
            ARCH="amd64"
            ;;
        arm64 | aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH_TYPE. Only amd64 and arm64 are supported."
            ;;
    esac

    info "Detected OS: $OS, Architecture: $ARCH"

    # 2. Get the latest version from GitHub Releases
    # Note: This requires the repository to have public releases.
    REPO="fun7257/sgv"
    LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | cut -d '"' -f 4)
    LATEST_VERSION=${LATEST_TAG#v}

    if [ -z "$LATEST_TAG" ]; then
        error "Could not fetch the latest version tag from GitHub. Please check the repository path and release status."
    fi
    
    info "Latest version is $LATEST_TAG"

    # 3. Download the pre-compiled binary
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/sgv_${LATEST_VERSION}_${OS}_${ARCH}.tar.gz"
    
    info "Downloading from $DOWNLOAD_URL..."
    TEMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TEMP_DIR"' EXIT # Cleanup temp directory on exit

    if ! curl -sSL "$DOWNLOAD_URL" -o "$TEMP_DIR/sgv.tar.gz"; then
        error "Download failed. Please check the URL and your network connection."
    fi

    # 4. Unpack and install
    info "Installing sgv to $INSTALL_DIR..."
    tar -xzf "$TEMP_DIR/sgv.tar.gz" -C "$TEMP_DIR"
    
    if [ ! -f "$TEMP_DIR/sgv" ]; then
        error "The downloaded archive does not contain the 'sgv' executable."
    fi

    # Move to install directory (requires sudo)
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TEMP_DIR/sgv" "$INSTALL_DIR/sgv"
    else
        info "Sudo privileges are required to install to $INSTALL_DIR."
        sudo mv "$TEMP_DIR/sgv" "$INSTALL_DIR/sgv"
    fi

    # 5. Set permissions and handle macOS Gatekeeper
    chmod +x "$INSTALL_DIR/sgv"
    info "Set executable permission for sgv."

    if [ "$OS" == "darwin" ]; then
        info "macOS detected. Removing quarantine attribute to prevent Gatekeeper issues."
        if xattr -d com.apple.quarantine "$INSTALL_DIR/sgv" 2>/dev/null; then
            info "Successfully removed quarantine attribute."
        else
            warn "Could not remove quarantine attribute. You might need to grant permissions manually in System Settings."
        fi
    fi

    # 6. Update Shell Configuration
    info "Adding environment variables to shell profile..."
    
    SHELL_CONFIG_FILE=""
    CURRENT_SHELL=$(basename "$SHELL")

    if [ "$CURRENT_SHELL" = "bash" ]; then
        SHELL_CONFIG_FILE="$HOME/.bashrc"
    elif [ "$CURRENT_SHELL" = "zsh" ]; then
        SHELL_CONFIG_FILE="$HOME/.zshrc"
    else
        warn "Could not detect a supported shell (bash or zsh). You will need to add the environment variables manually."
        warn "Add the following lines to your shell's startup file:"
        echo -e "\n# sgv (Simple Go Version) configuration\nexport GOROOT=\"$HOME/.sgv/current\"\nunset GOPATH\nexport PATH=\"$GOROOT/bin:$HOME/go/bin:$PATH\""
        return
    fi

    # Check if variables are already set
    if ! grep -q "export GOROOT" "$SHELL_CONFIG_FILE"; then
        echo -e "\n# sgv (Simple Go Version) configuration" >> "$SHELL_CONFIG_FILE"
        echo "export GOROOT=\"$HOME/.sgv/current\"" >> "$SHELL_CONFIG_FILE"
        echo "unset GOPATH" >> "$SHELL_CONFIG_FILE"
        echo "export PATH=\"$GOROOT/bin:$HOME/go/bin:$PATH\"" >> "$SHELL_CONFIG_FILE"
        info "Added GOROOT, unset GOPATH, and updated PATH in $SHELL_CONFIG_FILE."
    else
        info "Environment variables already seem to be set in $SHELL_CONFIG_FILE. Skipping."
    fi

    # --- Final Instructions ---
    echo -e "\n${GREEN}Installation successful!${NC}"
    echo -e "sgv has been installed to: ${YELLOW}$INSTALL_DIR/sgv${NC}"
    echo -e "\nPlease restart your terminal or run the following command to apply the changes:"
    echo -e "  ${YELLOW}source $SHELL_CONFIG_FILE${NC}"
    echo -e "\nThen you can start using sgv. For example:"
    echo -e "  ${YELLOW}sgv list${NC}"
}

# Run the main function
main
