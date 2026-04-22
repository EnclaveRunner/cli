package user

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <username>",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		RunE:  runDelete,
	}
}

func runDelete(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.UserColumns,
		os.Stdout,
	)

	u, err := c.DeleteUser(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return printer.Print([]any{u})
}
