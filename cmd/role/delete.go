package role

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <role>",
		Short: "Delete a role",
		Args:  cobra.ExactArgs(1),
		RunE:  runDelete,
	}
}

func runDelete(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.RoleColumns,
		os.Stdout,
	)

	r, err := c.DeleteRole(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	return printer.Print([]any{r})
}
