#!/usr/bin/env bash
set -e

# Enclave CLI Installation Script
# Usage: sh <(curl -L https://raw.githubusercontent.com/EnclaveRunner/CLI/main/install.sh)

GITHUB_REPO="EnclaveRunner/CLI"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="encl"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}==>${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

log_error() {
    echo -e "${RED}Error:${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        CYGWIN*|MINGW*|MSYS*) echo "windows";;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64";;
        aarch64|arm64)  echo "arm64";;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get latest release tag from GitHub
get_latest_release() {
    if command_exists curl; then
        curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
            grep '"tag_name":' | \
            sed -E 's/.*"([^"]+)".*/\1/'
    else
        log_error "curl is required but not installed"
        exit 1
    fi
}

# Download file with progress
download_file() {
    local url=$1
    local output=$2

    if command_exists curl; then
        echo "Downloading from $url..."
        curl -fsSL --proto '=https' --tlsv1.2 -o "$output" "$url"
    elif command_exists wget; then
        wget -q -O "$output" "$url"
    else
        log_error "curl or wget is required but neither is installed"
        exit 1
    fi
}



# Main installation function
install_cli() {
    log_info "Starting Enclave CLI installation..."

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    log_info "Detected system: ${OS}-${ARCH}"

    # Get latest release
    log_info "Fetching latest release information..."
    RELEASE_TAG=$(get_latest_release)
    if [ -z "$RELEASE_TAG" ]; then
        log_error "Failed to fetch latest release information"
        exit 1
    fi
    log_info "Latest release: ${RELEASE_TAG}"

    # Construct download URLs
    BINARY_FILENAME="enclave-cli-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_FILENAME="${BINARY_FILENAME}.exe"
    fi

    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${RELEASE_TAG}/${BINARY_FILENAME}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download binary
    log_info "Downloading ${BINARY_FILENAME}..."
    download_file "$DOWNLOAD_URL" "$TMP_DIR/$BINARY_FILENAME"

    # Make binary executable
    chmod +x "$TMP_DIR/$BINARY_FILENAME"

    # Install binary
    log_info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    mkdir -p "$INSTALL_DIR"
    mv "$TMP_DIR/$BINARY_FILENAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    log_info "Installation complete!"

    # Verify installation
    if command_exists "$BINARY_NAME"; then
        log_info "Enclave CLI version: $($BINARY_NAME --version 2>/dev/null || echo 'installed')"
    else
        log_warn "$INSTALL_DIR may not be in your PATH"
        log_warn "Add it to your PATH or use the full path: $INSTALL_DIR/$BINARY_NAME"
    fi

    echo ""
    log_info "Run '${BINARY_NAME} --help' to get started"
}

# Run installation
install_cli
