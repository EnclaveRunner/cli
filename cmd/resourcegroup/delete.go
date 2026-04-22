package resourcegroup

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a resource group",
		Args:  cobra.ExactArgs(1),
		RunE:  runDelete,
	}
}

func runDelete(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ResourceGroupColumns, os.Stdout)

	rg, err := c.DeleteResourceGroup(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("delete resource group: %w", err)
	}
	return printer.Print([]any{rg})
}
