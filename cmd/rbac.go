package cmd

import (
	"cli/client"
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rbacCmd = &cobra.Command{
	Use:   "rbac",
	Short: "Manage RBAC (Role-Based Access Control)",
	Long:  `Commands for managing roles, resource groups, policies, and user assignments.`,
}

// Role commands
var rbacRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage roles",
	Long:  `Commands for managing roles in the RBAC system.`,
}

var rbacRoleCreateCmd = &cobra.Command{
	Use:   "create <role>",
	Short: "Create a new role",
	Long:  `Create a new role in the system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		role := args[0]

		c := getClient()
		ctx := context.Background()

		body := client.RBACRole{
			Role: role,
		}

		resp, err := c.PostRbacRoleWithResponse(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create role")
			os.Exit(1)
		}

		ok := handleResponse(resp.HTTPResponse, "Role created successfully")
		if !ok {
			os.Exit(1)
		}
	},
}

var rbacRoleDeleteCmd = &cobra.Command{
	Use:   "delete <role>",
	Short: "Delete a role",
	Long:  `Delete a role from the system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		role := args[0]

		c := getClient()
		ctx := context.Background()

		body := client.RBACRole{
			Role: role,
		}

		resp, err := c.DeleteRbacRoleWithResponse(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete role")
			os.Exit(1)
		}

		ok := handleResponse(resp.HTTPResponse, "Role deleted successfully")
		if !ok {
			os.Exit(1)
		}
	},
}

var rbacRoleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles",
	Long:  `Retrieve a list of all roles in the system.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		ctx := context.Background()

		resp, err := c.GetRbacListRolesWithResponse(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list roles")
			os.Exit(1)
		}

		ok := handleResponse(resp.HTTPResponse, "")
		if !ok {
			os.Exit(1)
		} else {
			printSlice(*resp.JSON200)
		}
	},
}

var rbacRoleGetCmd = &cobra.Command{
	Use:   "get <role>",
	Short: "Get users assigned to a role",
	Long:  `Retrieve a list of users assigned to a specific role.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		role := args[0]

		c := getClient()
		ctx := context.Background()

		params := &client.GetRbacRoleParams{
			Role: role,
		}

		resp, err := c.GetRbacRoleWithResponse(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get role users")
			os.Exit(1)
		}

		ok := handleResponse(resp.HTTPResponse, "")
		if !ok {
			os.Exit(1)
		}

		users, err := getUsersByIds(ctx, *resp.JSON200)
		
		printUsers(users)
	},
}

// User role assignment commands
var rbacUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user role assignments",
	Long:  `Commands for managing role assignments to users.`,
}

var rbacUserAssignCmd = &cobra.Command{
	Use:   "assign <user-id> <role>",
	Short: "Assign a role to a user",
	Long:  `Assign a role to a specific user.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		userId := args[0]
		role := args[1]

		c := getClient()
		ctx := context.Background()

		body := client.PostRbacUserJSONRequestBody{
			UserId: userId,
			Role:   role,
		}

		resp, err := c.PostRbacUser(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to assign role to user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Role assigned to user successfully")
	},
}

var rbacUserRemoveCmd = &cobra.Command{
	Use:   "remove <user-id> <role>",
	Short: "Remove a role from a user",
	Long:  `Remove a role from a specific user.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		userId := args[0]
		role := args[1]

		c := getClient()
		ctx := context.Background()

		body := client.DeleteRbacUserJSONRequestBody{
			UserId: userId,
			Role:   role,
		}

		resp, err := c.DeleteRbacUser(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove role from user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Role removed from user successfully")
	},
}

var rbacUserGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get roles assigned to a user",
	Long:  `Retrieve a list of roles assigned to a specific user.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userId := args[0]

		c := getClient()
		ctx := context.Background()

		params := &client.GetRbacUserParams{
			UserId: userId,
		}

		resp, err := c.GetRbacUser(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user roles")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
	},
}

// Resource group commands
var rbacResourceGroupCmd = &cobra.Command{
	Use:   "resource-group",
	Short: "Manage resource groups",
	Long:  `Commands for managing resource groups in the RBAC system.`,
}

var rbacResourceGroupCreateCmd = &cobra.Command{
	Use:   "create <resource-group>",
	Short: "Create a new resource group",
	Long:  `Create a new resource group in the system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup := args[0]

		c := getClient()
		ctx := context.Background()

		body := client.PostRbacResourceGroupJSONRequestBody{
			ResourceGroup: resourceGroup,
		}

		resp, err := c.PostRbacResourceGroup(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create resource group")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Resource group created successfully")
	},
}

var rbacResourceGroupDeleteCmd = &cobra.Command{
	Use:   "delete <resource-group>",
	Short: "Delete a resource group",
	Long:  `Delete a resource group from the system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup := args[0]

		c := getClient()
		ctx := context.Background()

		body := client.DeleteRbacResourceGroupJSONRequestBody{
			ResourceGroup: resourceGroup,
		}

		resp, err := c.DeleteRbacResourceGroup(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete resource group")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Resource group deleted successfully")
	},
}

var rbacResourceGroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resource groups",
	Long:  `Retrieve a list of all resource groups in the system.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		ctx := context.Background()

		resp, err := c.GetRbacListResourceGroups(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list resource groups")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
	},
}

var rbacResourceGroupGetCmd = &cobra.Command{
	Use:   "get <resource-group>",
	Short: "Get endpoints in a resource group",
	Long:  `Retrieve a list of endpoints assigned to a specific resource group.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup := args[0]

		c := getClient()
		ctx := context.Background()

		params := &client.GetRbacResourceGroupParams{
			ResourceGroup: resourceGroup,
		}

		resp, err := c.GetRbacResourceGroup(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get resource group endpoints")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
	},
}

// Endpoint commands
var rbacEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Manage endpoint assignments",
	Long:  `Commands for managing endpoint assignments to resource groups.`,
}

var rbacEndpointAssignCmd = &cobra.Command{
	Use:   "assign <endpoint> <resource-group>",
	Short: "Assign an endpoint to a resource group",
	Long:  `Assign an endpoint to a specific resource group.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]
		resourceGroup := args[1]

		c := getClient()
		ctx := context.Background()

		body := client.PostRbacEndpointJSONRequestBody{
			Endpoint:      endpoint,
			ResourceGroup: resourceGroup,
		}

		resp, err := c.PostRbacEndpoint(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to assign endpoint to resource group")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Endpoint assigned to resource group successfully")
	},
}

var rbacEndpointRemoveCmd = &cobra.Command{
	Use:   "remove <endpoint> <resource-group>",
	Short: "Remove an endpoint from a resource group",
	Long:  `Remove an endpoint from its assigned resource group.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]
		resourceGroup := args[1]

		c := getClient()
		ctx := context.Background()

		body := client.DeleteRbacEndpointJSONRequestBody{
			Endpoint:      endpoint,
			ResourceGroup: resourceGroup,
		}

		resp, err := c.DeleteRbacEndpoint(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove endpoint from resource group")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Endpoint removed from resource group successfully")
	},
}

var rbacEndpointGetCmd = &cobra.Command{
	Use:   "get <endpoint>",
	Short: "Get resource group for an endpoint",
	Long:  `Retrieve the resource group assigned to a specific endpoint.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]

		c := getClient()
		ctx := context.Background()

		params := &client.GetRbacEndpointParams{
			Endpoint: endpoint,
		}

		resp, err := c.GetRbacEndpoint(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get endpoint resource group")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
	},
}

// Policy commands
var rbacPolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage RBAC policies",
	Long:  `Commands for managing RBAC policies that define permissions.`,
}

var rbacPolicyCreateCmd = &cobra.Command{
	Use:   "create <role> <resource-group> <permission>",
	Short: "Create a new RBAC policy",
	Long:  `Create a new RBAC policy that grants permissions to a role for a resource group.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		role := args[0]
		resourceGroup := args[1]
		permission := args[2]

		c := getClient()
		ctx := context.Background()

		body := client.RBACPolicy{
			Role:          role,
			ResourceGroup: resourceGroup,
			Permission:    client.RBACPolicyPermission(permission),
		}

		resp, err := c.PostRbacPolicy(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create policy")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Policy created successfully")
	},
}

var rbacPolicyDeleteCmd = &cobra.Command{
	Use:   "delete <role> <resource-group> <permission>",
	Short: "Delete an RBAC policy",
	Long:  `Delete an existing RBAC policy.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		role := args[0]
		resourceGroup := args[1]
		permission := args[2]

		c := getClient()
		ctx := context.Background()

		body := client.RBACPolicy{
			Role:          role,
			ResourceGroup: resourceGroup,
			Permission:    client.RBACPolicyPermission(permission),
		}

		resp, err := c.DeleteRbacPolicy(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete policy")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Policy deleted successfully")
	},
}

var rbacPolicyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all RBAC policies",
	Long:  `Retrieve a list of all RBAC policies in the system.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		ctx := context.Background()

		resp, err := c.GetRbacPolicy(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list policies")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
	},
}

func init() {
	rootCmd.AddCommand(rbacCmd)

	// Role commands
	rbacCmd.AddCommand(rbacRoleCmd)
	rbacRoleCmd.AddCommand(rbacRoleCreateCmd)
	rbacRoleCmd.AddCommand(rbacRoleDeleteCmd)
	rbacRoleCmd.AddCommand(rbacRoleListCmd)
	rbacRoleCmd.AddCommand(rbacRoleGetCmd)

	// User commands
	rbacCmd.AddCommand(rbacUserCmd)
	rbacUserCmd.AddCommand(rbacUserAssignCmd)
	rbacUserCmd.AddCommand(rbacUserRemoveCmd)
	rbacUserCmd.AddCommand(rbacUserGetCmd)

	// Resource group commands
	rbacCmd.AddCommand(rbacResourceGroupCmd)
	rbacResourceGroupCmd.AddCommand(rbacResourceGroupCreateCmd)
	rbacResourceGroupCmd.AddCommand(rbacResourceGroupDeleteCmd)
	rbacResourceGroupCmd.AddCommand(rbacResourceGroupListCmd)
	rbacResourceGroupCmd.AddCommand(rbacResourceGroupGetCmd)

	// Endpoint commands
	rbacCmd.AddCommand(rbacEndpointCmd)
	rbacEndpointCmd.AddCommand(rbacEndpointAssignCmd)
	rbacEndpointCmd.AddCommand(rbacEndpointRemoveCmd)
	rbacEndpointCmd.AddCommand(rbacEndpointGetCmd)

	// Policy commands
	rbacCmd.AddCommand(rbacPolicyCmd)
	rbacPolicyCmd.AddCommand(rbacPolicyCreateCmd)
	rbacPolicyCmd.AddCommand(rbacPolicyDeleteCmd)
	rbacPolicyCmd.AddCommand(rbacPolicyListCmd)
}
