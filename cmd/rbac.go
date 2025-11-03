package cmd

import (
	"cli/client"
	"context"
	"fmt"
	"os"

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

		successMsg := fmt.Sprintf("Role '%s' created", TextHighlight.Render(role))

		ok := handleResponse(resp, err, successMsg)
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

		successMsg := fmt.Sprintf("Role '%s deleted", TextHighlight.Render(role))

		ok := handleResponse(resp, err, successMsg)
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

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			roleInfo := getRoleInfo(ctx, *resp.JSON200)
			printRoles(roleInfo)
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

		ok := handleResponse(resp, err, "")
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
	Use:   "assign <username> <role>",
	Short: "Assign a role to a user",
	Long:  `Assign a role to a specific user.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		role := args[1]

		c := getClient()
		ctx := context.Background()

		user := getUserByName(ctx, username)

		assignReq, err := c.PostRbacUserWithResponse(ctx, client.PostRbacUserJSONRequestBody{
			UserId: user.Id,
			Role:   role,
		})

		successMsg := fmt.Sprintf("%s (%s) has role %s now", TextHighlight.Render(user.DisplayName), TextHighlight.Render(user.Id), TextHighlight.Render(role))

		ok := handleResponse(assignReq, err, successMsg)
		if !ok {
			os.Exit(1)
		}
	},
}

var rbacUserRemoveCmd = &cobra.Command{
	Use:   "remove <username> <role>",
	Short: "Remove a role from a user",
	Long:  `Remove a role from a specific user.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		role := args[1]

		c := getClient()
		ctx := context.Background()

		user := getUserByName(ctx, username)

		body := client.DeleteRbacUserJSONRequestBody{
			UserId: user.Id,
			Role:   role,
		}

		resp, err := c.DeleteRbacUserWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Role %s removed from %s (%s)", TextHighlight.Render(role), TextHighlight.Render(user.DisplayName), TextHighlight.Render(user.Id))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
	},
}

var rbacUserGetCmd = &cobra.Command{
	Use:   "get <username>",
	Short: "Get roles assigned to a user",
	Long:  `Retrieve a list of roles assigned to a specific user.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		c := getClient()
		ctx := context.Background()

		user := getUserByName(ctx, username)

		params := &client.GetRbacUserParams{
			UserId: user.Id,
		}

		resp, err := c.GetRbacUserWithResponse(ctx, params)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			roleInfo := getRoleInfo(ctx, *resp.JSON200)
			printRoles(roleInfo)
		}
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

		resp, err := c.PostRbacResourceGroupWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Resource group %s created", TextHighlight.Render(resourceGroup))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
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

		resp, err := c.DeleteRbacResourceGroupWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Resource group %s deleted", TextHighlight.Render(resourceGroup))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
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

		resp, err := c.GetRbacListResourceGroupsWithResponse(ctx)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			rgInfo := getResourceGroupInfo(ctx, *resp.JSON200)
			printResourceGroups(rgInfo)
		}
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

		resp, err := c.GetRbacResourceGroupWithResponse(ctx, params)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			printStringTable(*resp.JSON200, "ENDPOINT")
		}
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

		resp, err := c.PostRbacEndpointWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Endpoint %s assigned to resource group %s", TextHighlight.Render(endpoint), TextHighlight.Render(resourceGroup))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
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

		resp, err := c.DeleteRbacEndpointWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Endpoint %s removed from resource group %s", TextHighlight.Render(endpoint), TextHighlight.Render(resourceGroup))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
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

		resp, err := c.GetRbacEndpointWithResponse(ctx, params)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			rgInfo := getResourceGroupInfo(ctx, *resp.JSON200)
			printResourceGroups(rgInfo)
		}
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

		resp, err := c.PostRbacPolicyWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Policy created: role %s has %s permission on resource group %s", TextHighlight.Render(role), TextHighlight.Render(permission), TextHighlight.Render(resourceGroup))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
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

		resp, err := c.DeleteRbacPolicyWithResponse(ctx, body)

		successMsg := fmt.Sprintf("Policy deleted: role %s no longer has %s permission on resource group %s", TextHighlight.Render(role), TextHighlight.Render(permission), TextHighlight.Render(resourceGroup))

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
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

		resp, err := c.GetRbacPolicyWithResponse(ctx)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			printPolicies(*resp.JSON200)
		}
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
