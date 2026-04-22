package role

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <role>",
		Short: "Get a role by name",
		Args:  cobra.ExactArgs(1),
		RunE:  runGet,
	}
}

func runGet(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.RoleColumns,
		os.Stdout,
	)

	r, err := c.GetRole(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("get role: %w", err)
	}

	return printer.Print([]any{r})
}
