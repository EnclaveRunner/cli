# <img alt="Logo" width="80px" src="https://github.com/EnclaveRunner/.github/raw/main/img/enclave-logo.png" style="vertical-align: middle;" /> Enclave CLI
> [!WARNING]
> The enclave project is still under heavy development and object to changes. This can include APIs, schemas, interfaces and more. Productive usage is therefore not recommended yet (as long as no stable version is released).


## Quick Install

To install the latest version of Enclave CLI, run:

```bash
sh <(curl -L https://raw.githubusercontent.com/EnclaveRunner/cli/main/install.sh)
```

Or with wget:

```bash
sh <(wget -qO- https://raw.githubusercontent.com/EnclaveRunner/cli/main/install.sh)
```

## Custom Installation Directory

By default, the CLI is installed to `~/.local/bin/encl`. To install to a different directory:

```bash
INSTALL_DIR="/usr/local/bin" sh <(curl -L https://raw.githubusercontent.com/EnclaveRunner/cli/main/install.sh)
```

Make sure the installation directory is in your `PATH`.

## Manual Installation

### 1. Download the Binary

Go to the [releases page](https://github.com/EnclaveRunner/cli/releases) and download the appropriate binary for your system:

- **Linux AMD64**: `enclave-cli-linux-amd64`
- **Linux ARM64**: `enclave-cli-linux-arm64`
- **macOS AMD64**: `enclave-cli-darwin-amd64`
- **macOS ARM64** (Apple Silicon): `enclave-cli-darwin-arm64`

### 2. Install the Binary

```bash
# Make it executable
chmod +x enclave-cli-*

# Move to a directory in your PATH (recommended default)
mkdir -p ~/.local/bin
mv enclave-cli-* ~/.local/bin/encl

# Or to system-wide location (requires sudo)
sudo mv enclave-cli-* /usr/local/bin/encl
```

### 3. Verify Installation

```bash
encl --version
```

## Requirements

- curl or wget
- Linux, macOS, or Windows (WSL)
- AMD64 or ARM64 architecture

## Updating

To update to the latest version, simply run the installation script again. It will automatically download and install the latest release.

## Uninstalling

```bash
# If installed to ~/.local/bin (default)
rm ~/.local/bin/encl

# If installed to /usr/local/bin
sudo rm /usr/local/bin/encl
```

## Troubleshooting

### Command not found after installation

Make sure the installation directory is in your `PATH`:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="$HOME/.local/bin:$PATH"
```

Then reload your shell configuration:

```bash
source ~/.bashrc  # or source ~/.zshrc
```
