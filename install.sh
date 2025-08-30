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

# Function to update SGV shell configuration by replacing the entire configuration block
update_sgv_config() {
    local config_file="$1"

    info "Updating SGV shell configuration in $config_file..."

    # Remove all existing SGV configurations to avoid duplicates
    if grep -q ">>> SGV CONFIGURATION START <<<" "$config_file"; then
        info "Found existing SGV configuration. Removing it before adding the new one..."
        local temp_file=$(mktemp)
        sed '/# >>> SGV CONFIGURATION START <<</,/# >>> SGV CONFIGURATION END <<</d' "$config_file" > "$temp_file"
        mv "$temp_file" "$config_file"
        info "Removed old SGV configuration block."
    fi

    # Add new configuration at the end
    info "Adding new SGV configuration..."
    echo -e "\n# >>> SGV CONFIGURATION START <<<" >> "$config_file"
    echo "# sgv (Simple Go Version) configuration" >> "$config_file"
    echo "export GOROOT=\"\$HOME/.sgv/current\"" >> "$config_file"
    echo "export PATH=\"\$GOROOT/bin:\$HOME/go/bin:\$PATH\"" >> "$config_file"
    echo "unset GOPATH" >> "$config_file"
    echo "" >> "$config_file"
    echo "# SGV wrapper function for seamless environment variable loading" >> "$config_file"
    echo "sgv() {" >> "$config_file"
    echo "    command sgv \"\$@\"" >> "$config_file"
    echo "    local exit_code=\$?" >> "$config_file"
    echo "    # Auto-load environment variables after successful operations" >> "$config_file"
    echo "    if [ \$exit_code -eq 0 ]; then" >> "$config_file"
    echo "        # Check for version switch (direct version argument)" >> "$config_file"
    echo "        if [ \$# -eq 1 ] && [[ \"\$1\" =~ ^(go)?[0-9]+\\.[0-9]+(\\.[0-9]+)?\$ ]]; then" >> "$config_file"
    echo "            eval \"\$(command sgv env --shell --clean 2>/dev/null || true)\"" >> "$config_file"
    echo "        # Check for env command with write or unset flags" >> "$config_file"
    echo "        elif [ \"\$1\" = \"env\" ] && { [ \"\$2\" = \"-w\" ] || [ \"\$2\" = \"--write\" ] || [ \"\$2\" = \"-u\" ] || [ \"\$2\" = \"--unset\" ]; }; then" >> "$config_file"
    echo "            eval \"\$(command sgv env --shell 2>/dev/null || true)\"" >> "$config_file"
    echo "        # Check for auto, latest, and sub commands that may switch versions" >> "$config_file"
    echo "        elif [ \"\$1\" = \"auto\" ] || [ \"\$1\" = \"latest\" ] || [ \"\$1\" = \"sub\" ]; then" >> "$config_file"
    echo "            eval \"\$(command sgv env --shell --clean 2>/dev/null || true)\"" >> "$config_file"
    echo "        fi" >> "$config_file"
    echo "    fi" >> "$config_file"
    echo "    return \$exit_code" >> "$config_file"
    echo "}" >> "$config_file"
    echo "" >> "$config_file"
    echo "# Load SGV environment variables for current session" >> "$config_file"
    echo "if command -v sgv >/dev/null 2>&1 && [ -L \"\$HOME/.sgv/current\" ]; then" >> "$config_file"
    echo "    eval \"\$(command sgv env --shell --clean 2>/dev/null || true)\"" >> "$config_file"
    echo "fi" >> "$config_file"
    echo "# >>> SGV CONFIGURATION END <<<" >> "$config_file"
    info "Successfully added SGV configuration to $config_file."
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

    # 5. Set permissions
    chmod +x "$INSTALL_DIR/sgv"
    info "Set executable permission for sgv."

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
        echo -e "\n# sgv (Simple Go Version) configuration\nexport GOROOT=\"\$HOME/.sgv/current\"\nunset GOPATH\nexport PATH=\"\$GOROOT/bin:\$HOME/go/bin:\$PATH\""
        return
    fi

    # Update Shell Configuration
    update_sgv_config "$SHELL_CONFIG_FILE"

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
