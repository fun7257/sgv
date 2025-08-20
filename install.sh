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

# Function to clean up duplicate SGV configurations
clean_duplicate_sgv_config() {
    local config_file="$1"
    
    if [ ! -f "$config_file" ]; then
        return
    fi
    
    # Count how many sgv configuration markers exist
    local count=$(grep -c "sgv (Simple Go Version) configuration\|SGV wrapper function" "$config_file" 2>/dev/null || echo "0")
    
    if [ "$count" -gt 1 ]; then
        warn "Found multiple SGV configuration blocks in $config_file. Cleaning up duplicates..."
        
        # Create a backup
        local backup_file="$config_file.sgv-backup-$(date +%s)"
        cp "$config_file" "$backup_file"
        info "Created backup: $backup_file"
        
        # Remove all existing SGV configurations to avoid duplicates
        # We'll add the current configuration afterwards
        local temp_file=$(mktemp)
        local skip_lines=false
        
        while IFS= read -r line; do
            # Start skipping when we encounter SGV configuration
            if echo "$line" | grep -q "sgv (Simple Go Version) configuration\|SGV wrapper function"; then
                skip_lines=true
                continue
            fi
            
            # Continue skipping SGV-related lines
            if [ "$skip_lines" = true ]; then
                # Skip empty lines immediately after SGV marker
                if echo "$line" | grep -q "^\s*$"; then
                    continue
                fi
                
                # Skip SGV-related content
                if echo "$line" | grep -q "export GOROOT.*\.sgv\|export PATH.*\.sgv\|unset GOPATH\|sgv()\|command sgv\|eval.*sgv env\|Load SGV environment\|command -v sgv.*\.sgv\|return.*exit_code\|env --shell --clean\|env --shell\|\-\-clean"; then
                    continue
                fi
                
                # Skip control structures and braces that are part of the function
                if echo "$line" | grep -q "^\s*{\s*$\|^\s*}\s*$\|^\s*if \[\|^\s*elif \[\|^\s*fi\s*$\|local exit_code"; then
                    continue
                fi
                
                # If we hit a new section (not SGV related), stop skipping
                if echo "$line" | grep -q "^[[:space:]]*[^#[:space:]]" && ! echo "$line" | grep -qi "sgv"; then
                    skip_lines=false
                fi
            fi
            
            # Keep the line if we're not skipping
            if [ "$skip_lines" = false ]; then
                echo "$line" >> "$temp_file"
            fi
        done < "$config_file"
        
        # Replace the original file with cleaned version
        mv "$temp_file" "$config_file"
        info "Removed duplicate SGV configurations from $config_file"
    fi
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

    # Clean up any duplicate SGV configurations before adding new ones
    clean_duplicate_sgv_config "$SHELL_CONFIG_FILE"

    # Check if SGV configuration is already installed
    if ! grep -q "sgv (Simple Go Version) configuration" "$SHELL_CONFIG_FILE" && ! grep -q "SGV wrapper function" "$SHELL_CONFIG_FILE"; then
        echo -e "\n# sgv (Simple Go Version) configuration" >> "$SHELL_CONFIG_FILE"
        echo "export GOROOT=\"\$HOME/.sgv/current\"" >> "$SHELL_CONFIG_FILE"
        echo "export PATH=\"\$GOROOT/bin:\$HOME/go/bin:\$PATH\"" >> "$SHELL_CONFIG_FILE"
        echo "unset GOPATH" >> "$SHELL_CONFIG_FILE"
        echo "" >> "$SHELL_CONFIG_FILE"
        echo "# SGV wrapper function for seamless environment variable loading" >> "$SHELL_CONFIG_FILE"
        echo "sgv() {" >> "$SHELL_CONFIG_FILE"
        echo "    command sgv \"\$@\"" >> "$SHELL_CONFIG_FILE"
        echo "    local exit_code=\$?" >> "$SHELL_CONFIG_FILE"
        echo "    # Auto-load environment variables after successful operations" >> "$SHELL_CONFIG_FILE"
        echo "    if [ \$exit_code -eq 0 ]; then" >> "$SHELL_CONFIG_FILE"
        echo "        # Check for version switch (direct version argument)" >> "$SHELL_CONFIG_FILE"
        echo "        if [ \$# -eq 1 ] && [[ \"\$1\" =~ ^go[0-9]+\\.[0-9]+(\\.[0-9]+)?\$ ]]; then" >> "$SHELL_CONFIG_FILE"
        echo "            eval \"\$(command sgv env --shell --clean 2>/dev/null || true)\"" >> "$SHELL_CONFIG_FILE"
        echo "        # Check for env command with write or unset flags" >> "$SHELL_CONFIG_FILE"
        echo "        elif [ \"\$1\" = \"env\" ] && ([ \"\$2\" = \"-w\" ] || [ \"\$2\" = \"--write\" ] || [ \"\$2\" = \"-u\" ] || [ \"\$2\" = \"--unset\" ]); then" >> "$SHELL_CONFIG_FILE"
        echo "            eval \"\$(command sgv env --shell 2>/dev/null || true)\"" >> "$SHELL_CONFIG_FILE"
        echo "        # Check for auto and latest commands that may switch versions" >> "$SHELL_CONFIG_FILE"
        echo "        elif [ \"\$1\" = \"auto\" ] || [ \"\$1\" = \"latest\" ]; then" >> "$SHELL_CONFIG_FILE"
        echo "            eval \"\$(command sgv env --shell --clean 2>/dev/null || true)\"" >> "$SHELL_CONFIG_FILE"
        echo "        fi" >> "$SHELL_CONFIG_FILE"
        echo "    fi" >> "$SHELL_CONFIG_FILE"
        echo "    return \$exit_code" >> "$SHELL_CONFIG_FILE"
        echo "}" >> "$SHELL_CONFIG_FILE"
        echo "" >> "$SHELL_CONFIG_FILE"
        echo "# Load SGV environment variables for current session" >> "$SHELL_CONFIG_FILE"
        echo "if command -v sgv >/dev/null 2>&1 && [ -L \"\$HOME/.sgv/current\" ]; then" >> "$SHELL_CONFIG_FILE"
        echo "    eval \"\$(command sgv env --shell --clean 2>/dev/null || true)\"" >> "$SHELL_CONFIG_FILE"
        echo "fi" >> "$SHELL_CONFIG_FILE"
        info "Added GOROOT, unset GOPATH, updated PATH, and enabled seamless environment variable loading in $SHELL_CONFIG_FILE."
    else
        info "SGV configuration already exists in $SHELL_CONFIG_FILE. Skipping shell configuration."
        
        # Check if the existing configuration is outdated and suggest manual update
        if grep -q "SGV wrapper function" "$SHELL_CONFIG_FILE"; then
            # Check if the wrapper function includes the latest environment variable loading logic
            if ! grep -A 25 "SGV wrapper function" "$SHELL_CONFIG_FILE" | grep -q "env --shell --clean"; then
                warn "Your SGV shell configuration might be outdated."
                warn "Consider removing the old SGV configuration from $SHELL_CONFIG_FILE and re-running this installer."
            fi
            # Also check if the session loading part uses --clean
            if ! grep -A 5 "Load SGV environment" "$SHELL_CONFIG_FILE" | grep -q "env --shell --clean"; then
                warn "Your SGV session loading configuration might be outdated."
                warn "Consider removing the old SGV configuration from $SHELL_CONFIG_FILE and re-running this installer."
            fi
        fi
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
