# Enclave CLI Configuration

This package handles configuration management for the Enclave CLI using [Viper](https://github.com/spf13/viper).

## Configuration Sources

The CLI loads configuration from multiple sources with the following precedence (highest to lowest):

1. **CLI Flags** (highest precedence)
2. **Environment Variables**
3. **Config File in Current Directory** (`./cli.yaml`)
4. **Config File in User Home** (`$HOME/.enclave/cli.yaml`)
5. **Config File in System Directory** (`/etc/enclave/cli.yaml`)
6. **Default Values** (lowest precedence)

## Configuration Options

### API Server URL

The API server endpoint for the Enclave platform.

- **Config Key**: `api_server_url`
- **CLI Flag**: `--api-url`
- **Environment Variable**: `ENCLAVE_API_SERVER_URL`
- **Default**: `https://api.enclave.io`

### Authentication

Basic authentication credentials for API requests.

#### Username

- **Config Key**: `auth.username`
- **CLI Flag**: `--auth-username`
- **Environment Variable**: `ENCLAVE_AUTH_USERNAME`
- **Default**: None

#### Password

- **Config Key**: `auth.password`
- **CLI Flag**: `--auth-password`
- **Environment Variable**: `ENCLAVE_AUTH_PASSWORD`
- **Default**: None

## Usage Examples

### Using CLI Flags

```bash
enclave config --api-url https://api.example.com --auth-username myuser --auth-password mypass
```

### Using Environment Variables

```bash
export ENCLAVE_API_SERVER_URL=https://api.example.com
export ENCLAVE_AUTH_USERNAME=myuser
export ENCLAVE_AUTH_PASSWORD=mypass
enclave config
```

### Using a Config File

Create a file at one of the supported locations:

**Example: `./cli.yaml`**

```yaml
api_server_url: https://api.example.com

auth:
  username: myuser
  password: mypass
```

Then run:

```bash
enclave config
```

### Using a Custom Config File

```bash
enclave config --config /path/to/custom/config.yaml
```

### Combining Multiple Sources

You can combine different sources, with higher precedence sources overriding lower ones:

```bash
# Config file sets base configuration
cat > cli.yaml << EOF
api_server_url: https://api.example.com
auth:
  username: default-user
  password: default-pass
EOF

# Environment variable overrides username
export ENCLAVE_AUTH_USERNAME=env-user

# CLI flag overrides API URL
enclave config --api-url https://override.example.com

# Result:
# - API URL: https://override.example.com (from CLI flag)
# - Username: env-user (from environment variable)
# - Password: default-pass (from config file)
```

## Configuration File Format

The configuration file uses YAML format. Here's a complete example:

```yaml
# API Server Configuration
api_server_url: https://api.enclave.io

# Authentication Configuration
auth:
  username: your-username
  password: your-password
```

A template file is provided as `cli.yaml.example` in the root of the project.

## Accessing Configuration in Code

The global configuration is available via the `config.Cfg` variable after `config.Init()` has been called:

```go
import "cli/config"

// Access the API server URL
url := config.Cfg.APIServerURL

// Access authentication
if config.Cfg.Auth != nil {
    authHeader := config.Cfg.Auth.GetAuthHeader()
}
```

## Security Best Practices

1. **Don't commit credentials**: Never commit config files containing real credentials to version control
2. **Use environment variables**: For CI/CD and production environments, prefer environment variables over config files
3. **File permissions**: If using config files with credentials, ensure they have restrictive permissions:
   ```bash
   chmod 600 ~/.enclave/cli.yaml
   ```
4. **Sensitive data**: The `config` command hides password values when displaying configuration

## Initialization

The configuration system is automatically initialized in the `PersistentPreRunE` hook of the root command, ensuring it's available for all subcommands.

```go
func Init(viperInstance ...*viper.Viper) error
```

The `Init` function:
- Accepts an optional Viper instance for flag binding
- Sets up config file search paths
- Configures environment variable support
- Reads and parses the config file (if found)
- Unmarshals configuration into the `Cfg` global variable
- Returns an error if config parsing fails (only when using `--config` flag)

## Testing Configuration

Use the `config` subcommand to verify your configuration:

```bash
enclave config
```

This displays the current configuration loaded from all sources.