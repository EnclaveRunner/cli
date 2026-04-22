package user

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <username>",
		Short: "Get a user by username",
		Args:  cobra.ExactArgs(1),
		RunE:  runGet,
	}
}

func runGet(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.UserColumns, os.Stdout)

	u, err := c.GetUser(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	return printer.Print([]any{u})
}
