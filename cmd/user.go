package cmd

import (
	"bytes"
	"cli/client"
	"cli/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

		resp, err := c.PostUsersUser(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "User created successfully")
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete <user-id>",
	Short: "Delete a user",
	Long:  `Delete a user by their ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userId := args[0]

		c := getClient()
		ctx := context.Background()

		body := client.DeleteUsersUserJSONRequestBody{
			Id: userId,
		}

		resp, err := c.DeleteUsersUser(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "User deleted successfully")
	},
}

var userUpdateCmd = &cobra.Command{
	Use:   "update <user-id>",
	Short: "Update a user",
	Long:  `Update a user's name, display name, or password.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userId := args[0]
		newName, _ := cmd.Flags().GetString("new-name")
		newDisplayName, _ := cmd.Flags().GetString("new-display-name")
		newPassword, _ := cmd.Flags().GetString("new-password")

		if newName == "" && newDisplayName == "" && newPassword == "" {
			log.Error().Msg("at least one of --new-name, --new-display-name, or --new-password must be provided")
			os.Exit(1)
		}

		c := getClient()
		ctx := context.Background()

		body := client.PatchUser{
			Id: userId,
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

		resp, err := c.PatchUsersUser(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "User updated successfully")
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get user information",
	Long:  `Retrieve information about a specific user by their ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userId := args[0]

		c := getClient()
		ctx := context.Background()

		params := &client.GetUsersUserParams{
			UserId: userId,
		}

		resp, err := c.GetUsersUser(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
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

		resp, err := c.GetUsersList(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list users")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
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

		resp, err := c.GetUsersMe(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "")
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
			log.Error().Msg("at least one of --new-name, --new-display-name, or --new-password must be provided")
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

		resp, err := c.PatchUsersMe(ctx, body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update current user")
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, "Current user updated successfully")
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

func getClient() *client.Client {
	if config.Cfg.APIServerURL == "" {
		log.Error().Msg("API server URL not configured")
		os.Exit(1)
	}

	if config.Cfg.Auth == nil {
		log.Error().Msg("Authentication not configured")
		os.Exit(1)
	}

	c, err := client.NewClient(config.Cfg.APIServerURL, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", config.Cfg.Auth.GetAuthHeader())
		return nil
	}))

	if err != nil {
		log.Error().Err(err).Msg("Failed to create API client")
		os.Exit(1)
	}

	return c
}

func handleResponse(resp *http.Response, successMsg string) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		os.Exit(1)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if successMsg != "" {
			log.Info().Msg(successMsg)
		}
		if len(body) > 0 {
			// Try to pretty print JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
				fmt.Println(prettyJSON.String())
			} else {
				fmt.Println(string(body))
			}
		}
	} else {
		log.Error().Int("status", resp.StatusCode).Msg("Request failed")
		if len(body) > 0 {
			fmt.Fprintln(os.Stderr, string(body))
		}
		os.Exit(1)
	}
}
