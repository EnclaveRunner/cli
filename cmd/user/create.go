package user

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <username> <display-name> <password>",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(3),
		RunE:  runCreate,
	}
}

func runCreate(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.UserColumns, os.Stdout)

	u, err := c.CreateUser(cmd.Context(), args[0], args[2], args[1])
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return printer.Print([]any{u})
}
