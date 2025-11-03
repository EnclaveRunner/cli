package cmd

import (
	"cli/client"
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  `Commands for managing users in the Enclave system.`,
}

var userCreateCmd = &cobra.Command{
	Use:   "create <name> <display-name> <password>",
	Short: "Create a new user",
	Long:  `Create a new user with the specified name, display name, and password.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		displayName := args[1]
		password := args[2]

		c := getClient()
		ctx := context.Background()

		body := client.CreateUser{
			Name:        name,
			DisplayName: displayName,
			Password:    password,
		}

		resp, err := c.PostUsersUserWithResponse(ctx, body)

		successMsg := "User created successfully"

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete <username>",
	Short: "Delete a user",
	Long:  `Delete a user by their username.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		c := getClient()
		ctx := context.Background()

		user := getUserByName(ctx, username)

		body := client.DeleteUsersUserJSONRequestBody{
			Id: user.Id,
		}

		resp, err := c.DeleteUsersUserWithResponse(ctx, body)

		successMsg := "User deleted successfully"

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
	},
}

var userUpdateCmd = &cobra.Command{
	Use:   "update <username>",
	Short: "Update a user",
	Long:  `Update a user's name, display name, or password.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		newName, _ := cmd.Flags().GetString("new-name")
		newDisplayName, _ := cmd.Flags().GetString("new-display-name")
		newPassword, _ := cmd.Flags().GetString("new-password")

		if newName == "" && newDisplayName == "" && newPassword == "" {
			log.Error().
				Msg("at least one of --new-name, --new-display-name, or --new-password must be provided")
			os.Exit(1)
		}

		c := getClient()
		ctx := context.Background()

		user := getUserByName(ctx, username)

		body := client.PatchUser{
			Id: user.Id,
		}

		if newName != "" {
			body.NewName = &newName
		}
		if newDisplayName != "" {
			body.NewDisplayName = &newDisplayName
		}
		if newPassword != "" {
			body.NewPassword = &newPassword
		}

		resp, err := c.PatchUsersUserWithResponse(ctx, body)

		successMsg := "User updated successfully"

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get <username>",
	Short: "Get user information",
	Long:  `Retrieve information about a specific user by their username.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		ctx := context.Background()

		user := getUserByName(ctx, username)

		printUser(user)
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `Retrieve a list of all users in the system.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		ctx := context.Background()

		resp, err := c.GetUsersListWithResponse(ctx)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			// Convert []UserResponse to []*UserResponse
			users := make([]*client.UserResponse, len(*resp.JSON200))
			for i := range *resp.JSON200 {
				users[i] = &(*resp.JSON200)[i]
			}
			printUsers(users)
		}
	},
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Manage current user",
	Long:  `Commands for managing the currently authenticated user.`,
}

var userMeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current user information",
	Long:  `Retrieve information about the currently authenticated user.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		ctx := context.Background()

		resp, err := c.GetUsersMeWithResponse(ctx)

		ok := handleResponse(resp, err, "")
		if !ok {
			os.Exit(1)
		} else {
			printUser(resp.JSON200)
		}
	},
}

var userMeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update current user",
	Long:  `Update the currently authenticated user's name, display name, or password.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		newName, _ := cmd.Flags().GetString("new-name")
		newDisplayName, _ := cmd.Flags().GetString("new-display-name")
		newPassword, _ := cmd.Flags().GetString("new-password")

		if newName == "" && newDisplayName == "" && newPassword == "" {
			log.Error().
				Msg("at least one of --new-name, --new-display-name, or --new-password must be provided")
			os.Exit(1)
		}

		c := getClient()
		ctx := context.Background()

		body := client.PatchMe{}

		if newName != "" {
			body.NewName = &newName
		}
		if newDisplayName != "" {
			body.NewDisplayName = &newDisplayName
		}
		if newPassword != "" {
			body.NewPassword = &newPassword
		}

		resp, err := c.PatchUsersMeWithResponse(ctx, body)

		successMsg := "Current user updated successfully"

		ok := handleResponse(resp, err, successMsg)
		if !ok {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Create command
	userCmd.AddCommand(userCreateCmd)

	// Delete command
	userCmd.AddCommand(userDeleteCmd)

	// Update command
	userCmd.AddCommand(userUpdateCmd)
	userUpdateCmd.Flags().String("new-name", "", "New user name")
	userUpdateCmd.Flags().String("new-display-name", "", "New display name")
	userUpdateCmd.Flags().String("new-password", "", "New password")

	// Get command
	userCmd.AddCommand(userGetCmd)

	// List command
	userCmd.AddCommand(userListCmd)

	// Me command
	userCmd.AddCommand(userMeCmd)
	userMeCmd.AddCommand(userMeGetCmd)
	userMeCmd.AddCommand(userMeUpdateCmd)
	userMeUpdateCmd.Flags().String("new-name", "", "New user name")
	userMeUpdateCmd.Flags().String("new-display-name", "", "New display name")
	userMeUpdateCmd.Flags().String("new-password", "", "New password")
}
