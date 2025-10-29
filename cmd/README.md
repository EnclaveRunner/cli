# Enclave CLI Commands

This document provides examples of all available commands in the Enclave CLI.

## User Management

### Create a user
```bash
enclave user create <name> <display-name> <password>
```
Example:
```bash
enclave user create john "John Doe" secretpassword123
```

### Delete a user
```bash
enclave user delete <user-id>
```
Example:
```bash
enclave user delete 550e8400-e29b-41d4-a716-446655440000
```

### Update a user
```bash
enclave user update <user-id> [flags]
```
Flags:
- `--new-name` - New user name
- `--new-display-name` - New display name
- `--new-password` - New password

Example:
```bash
enclave user update 550e8400-e29b-41d4-a716-446655440000 --new-display-name "John Smith"
enclave user update 550e8400-e29b-41d4-a716-446655440000 --new-name johnsmith --new-password newpass123
```

### Get user information
```bash
enclave user get <user-id>
```
Example:
```bash
enclave user get 550e8400-e29b-41d4-a716-446655440000
```

### List all users
```bash
enclave user list
```

### Get current user information
```bash
enclave user me get
```

### Update current user
```bash
enclave user me update [flags]
```
Flags:
- `--new-name` - New user name
- `--new-display-name` - New display name
- `--new-password` - New password

Example:
```bash
enclave user me update --new-display-name "Jane Doe"
```

## RBAC - Role Management

### Create a role
```bash
enclave rbac role create <role>
```
Example:
```bash
enclave rbac role create admin
```

### Delete a role
```bash
enclave rbac role delete <role>
```
Example:
```bash
enclave rbac role delete admin
```

### List all roles
```bash
enclave rbac role list
```

### Get users assigned to a role
```bash
enclave rbac role get <role>
```
Example:
```bash
enclave rbac role get admin
```

## RBAC - User Role Assignments

### Assign a role to a user
```bash
enclave rbac user assign <user-id> <role>
```
Example:
```bash
enclave rbac user assign 550e8400-e29b-41d4-a716-446655440000 admin
```

### Remove a role from a user
```bash
enclave rbac user remove <user-id> <role>
```
Example:
```bash
enclave rbac user remove 550e8400-e29b-41d4-a716-446655440000 admin
```

### Get roles assigned to a user
```bash
enclave rbac user get <user-id>
```
Example:
```bash
enclave rbac user get 550e8400-e29b-41d4-a716-446655440000
```

## RBAC - Resource Group Management

### Create a resource group
```bash
enclave rbac resource-group create <resource-group>
```
Example:
```bash
enclave rbac resource-group create api-endpoints
```

### Delete a resource group
```bash
enclave rbac resource-group delete <resource-group>
```
Example:
```bash
enclave rbac resource-group delete api-endpoints
```

### List all resource groups
```bash
enclave rbac resource-group list
```

### Get endpoints in a resource group
```bash
enclave rbac resource-group get <resource-group>
```
Example:
```bash
enclave rbac resource-group get api-endpoints
```

## RBAC - Endpoint Management

### Assign an endpoint to a resource group
```bash
enclave rbac endpoint assign <endpoint> <resource-group>
```
Example:
```bash
enclave rbac endpoint assign /api/v1/users api-endpoints
```

### Remove an endpoint from a resource group
```bash
enclave rbac endpoint remove <endpoint> <resource-group>
```
Example:
```bash
enclave rbac endpoint remove /api/v1/users api-endpoints
```

### Get resource group for an endpoint
```bash
enclave rbac endpoint get <endpoint>
```
Example:
```bash
enclave rbac endpoint get /api/v1/users
```

## RBAC - Policy Management

### Create an RBAC policy
```bash
enclave rbac policy create <role> <resource-group> <permission>
```
Permissions: `GET`, `POST`, `PATCH`, `DELETE`, `HEAD`, or `*` (all permissions)
Use `*` for role or resource-group to match all.

Examples:
```bash
# Grant GET permission to admin role for api-endpoints resource group
enclave rbac policy create admin api-endpoints GET

# Grant all permissions to superadmin role for all resource groups
enclave rbac policy create superadmin "*" "*"

# Grant POST and PATCH permissions to editor role for api-endpoints
enclave rbac policy create editor api-endpoints POST
enclave rbac policy create editor api-endpoints PATCH
```

### Delete an RBAC policy
```bash
enclave rbac policy delete <role> <resource-group> <permission>
```
Example:
```bash
enclave rbac policy delete admin api-endpoints GET
```

### List all RBAC policies
```bash
enclave rbac policy list
```

## Configuration

The CLI requires configuration for the API server URL and authentication credentials. These can be provided via:

1. **Flags** (highest priority):
   ```bash
   enclave --api-url https://api.example.com --auth-username admin --auth-password pass user list
   ```

2. **Environment variables**:
   ```bash
   export ENCLAVE_API_SERVER_URL=https://api.example.com
   export ENCLAVE_AUTH_USERNAME=admin
   export ENCLAVE_AUTH_PASSWORD=pass
   enclave user list
   ```

3. **Configuration file** (lowest priority):
   Create a file at `./cli.yml`, `~/.enclave/cli.yml`, or `/etc/enclave/cli.yml`:
   ```yaml
   api_server_url: https://api.example.com
   auth:
     username: admin
     password: pass
   ```

### View current configuration
```bash
enclave config
```

### Use a custom config file
```bash
enclave --config /path/to/config.yml user list
```
