package role

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <role>",
		Short: "Create a new role",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreate,
	}
	cmd.Flags().StringSlice("users", nil, "Users to assign to the role")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.RoleColumns,
		os.Stdout,
	)

	users, _ := cmd.Flags().GetStringSlice("users")
	r, err := c.CreateRole(cmd.Context(), args[0], users)
	if err != nil {
		return fmt.Errorf("create role: %w", err)
	}

	return printer.Print([]any{r})
}
